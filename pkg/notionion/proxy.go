package notionion

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/elazarl/goproxy"
	"github.com/jomei/notionapi"
)

//ProxyRequestHTTPHandler: Proxy handler sending request to notion page
func ProxyRequestHTTPHandler(client *notionapi.Client, pageid string, codeReq notionapi.CodeBlock, codeResp notionapi.CodeBlock) (h goproxy.ReqHandler) {
	var f goproxy.FuncReqHandler
	f = func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if active, err := RequestProxyStatus(client, pageid); err != nil {
			fmt.Println(err)
			return r, nil
		} else if active {
			//reset response section
			ClearResponseCode(client, codeResp.ID)
			// Print request on Notion proxy page
			r.Header.Del("Content-Lenght")
			requestDump, err := httputil.DumpRequest(r, true) //Use dumprequestOut to have the full request + use string(requestDump)
			if err != nil {
				fmt.Println(err)
			}
			if req, err := getRequestWithoutContentLength(requestDump); err != nil {
				fmt.Println(err)
				return r, nil
			} else {
				UpdateCodeContent(client, codeReq.ID, req)
			}

			//UpdateCodeContent(client, codeReq.ID, string(requestDump))
			//+enable button
			if err := EnableRequestButtons(client, pageid); err != nil {
				fmt.Println(err)
			}
			//wait for action (forward or drop)
			action := WaitAction(client, pageid)

			//disable button
			if err := DisableRequestButtons(client, pageid); err != nil {
				fmt.Println(err)
			}

			switch action {
			case FORWARD:
				//todo: retrieve code content -> to string
				reqFromPage, err := RequestRequestCodeContent(client, pageid)
				if err != nil {
					fmt.Println("Failed to retrieve request from notion proxy page:", err)
				}
				reqFromPage, err = addContentLength([]byte(reqFromPage)) //comment it, if you don't use getRequestWithoutContentLength
				if err != nil {
					fmt.Println(err)
				}
				reader := bufio.NewReader(strings.NewReader(reqFromPage))
				if r, err = http.ReadRequest(reader); err != nil {
					fmt.Println("Failed parsing request from notion proxy page:", err)
				}
				return r, nil
			case DROP:
				return nil, nil
			}
		}
		return r, nil
	}
	return f
}

//ProxyResponseHTTPHandler: Proxy handler sending response to notion page
func ProxyResponseHTTPHandler(client *notionapi.Client, pageid string, codeResp notionapi.CodeBlock) (h goproxy.RespHandler) {
	var f goproxy.FuncRespHandler
	f = func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if active, err := RequestProxyStatus(client, pageid); err != nil {
			fmt.Println(err)
			return resp
		} else if active && resp != nil {
			// Print response on Notion proxy page
			responseDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				fmt.Println(err)
			}

			// Print response in Notion proxy page
			UpdateCodeContent(client, codeResp.ID, string(responseDump))
		}
		return resp
	}
	return f
}

//getRequestWithoutContentLength: withdraw content-lenght header from a request (from:https://github.com/ariary/HTTPCustomHouse/blob/main/pkg/parser/parser.go)
func getRequestWithoutContentLength(requestDump []byte) (req string, err error) {
	reader := bytes.NewReader(requestDump)

	tp := textproto.NewReader(bufio.NewReader(reader))
	// First line: POST /index.html HTTP/1.0 or other
	if s, err := tp.ReadLine(); err != nil {
		return "", err
	} else {
		req += s + "\r\n"
	}

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return "", err
	}
	// Get header + delete Content-Length
	httpHeader := http.Header(mimeHeader)
	delete(httpHeader, "Content-Length")
	//append headers
	// always append Host first even if it normally does not have significancy
	// append 1 of them, delete it and continue
	hosts := httpHeader["Host"]
	if len(hosts) != 0 {
		req += fmt.Sprintf("Host: %s\r\n", hosts[0])
	}
	if len(hosts) > 1 { //several Host headers
		httpHeader["Host"] = hosts[1:]
	}
	for h, values := range httpHeader { // append other  http header
		for i := 0; i < len(values); i++ {
			req += fmt.Sprintf("%s: %s\r\n", h, values[i])
		}
	}

	//Get body
	bodyB, err := io.ReadAll(tp.R)
	if err != nil {
		return "", err
	}
	bodyB = append([]byte("\r\n"), bodyB...)
	req += string(bodyB)

	return req, nil
}

//addContentLength: add content-lenght header from a request to make the body usable
func addContentLength(requestDump []byte) (req string, err error) {
	reader := bytes.NewReader(requestDump)

	tp := textproto.NewReader(bufio.NewReader(reader))
	// First line: POST /index.html HTTP/1.0 or other
	if first, err := tp.ReadLine(); err != nil {
		return "", err
	} else {
		req += first + "\n"
	}

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return "", err
	}
	// Get header + delete Content-Length
	httpHeader := http.Header(mimeHeader)

	// Get body
	bodyB, err := io.ReadAll(tp.R)
	if err != nil {
		return "", err
	}
	cl := len(string(bodyB))
	bodyB = append([]byte("\r\n"), bodyB...)
	body := string(bodyB)

	// Add content-length
	httpHeader.Add("Content-Length", strconv.Itoa(cl))

	// Reconstruct request
	//append headers
	hosts := httpHeader["Host"]
	if len(hosts) != 0 {
		req += fmt.Sprintf("Host: %s\r\n", hosts[0])
	}
	if len(hosts) > 1 { //several Host headers
		httpHeader["Host"] = hosts[1:]
	}
	for h, values := range httpHeader { // append other  http header
		for i := 0; i < len(values); i++ {
			req += fmt.Sprintf("%s: %s\r\n", h, values[i])
		}
	}
	//append body
	req += body

	return req, nil
}

//WaitAction: If we enter this function, the request is in the notion page and the user
// is treated it. We are waiting for the user check neither FORWARD or DROP.
func WaitAction(client *notionapi.Client, pageid string) string {
	//spinner config
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond) // or spinner.CharSets[39]
	s.Color("blue")
	s.Suffix = " Request is being treated in Notion proxy page"

	//channel config
	stopchan := make(chan struct{}) // a channel to tell it to stop
	actionChoice := make(chan string)

	s.Start()

	go ListenDropButton(client, pageid, stopchan, actionChoice)
	go ListenForwardButton(client, pageid, stopchan, actionChoice)

	<-stopchan // wait for it to have stopped
	action := <-actionChoice
	s.FinalMSG = "✔ Request treated: " + action + "\n"
	s.Stop()
	return FORWARD
}

//ListenForwardButton: function that constantly check the forward button status to see if it is checked. Update channels consequently
func ListenForwardButton(client *notionapi.Client, pageid string, stopchan chan struct{}, action chan<- string) {
	// defer close(stoppedchan)
	// defer func() {

	// }()
	for {
		select {
		default:
			if check, err := RequestForwardButtonStatus(client, pageid); err != nil {
				fmt.Println(err)
			} else if check {
				close(stopchan)
				action <- FORWARD
			} else {
				//time between each request
				time.Sleep(2 * time.Second) //todo: customizable
			}
		case <-stopchan:
			// stop
			return
		}
	}
}

//ListenDropButton: function that constantly check the drop button status to see if it is checked. Update channels consequently
func ListenDropButton(client *notionapi.Client, pageid string, stopchan chan struct{}, action chan<- string) {
	// defer close(stoppedchan)
	// defer func() {

	// }()
	for {
		select {
		default:
			if check, err := RequestDropButtonStatus(client, pageid); err != nil {
				fmt.Println(err)
			} else if check {
				close(stopchan)
				action <- DROP
			} else {
				//time between each request
				time.Sleep(2 * time.Second) //todo: customizable
			}
		case <-stopchan:
			// stop
			return
		}
	}
}
