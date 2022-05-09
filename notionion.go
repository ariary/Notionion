package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/elazarl/goproxy"
	"github.com/jomei/notionapi"
)

func main() {
	// integration token
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("❌ Please set NOTION_TOKEN envvar with your integration token before launching notionion")
		os.Exit(92)
	}
	// page id
	pageurl := os.Getenv("NOTION_PAGE_URL")
	if pageurl == "" {
		fmt.Println("❌ Please set NOTION_PAGE_URL envvar with your page id before launching notionion (CTRL+L on desktop app)")
		os.Exit(92)
	}

	pageid := pageurl[strings.LastIndex(pageurl, "-")+1:]
	if pageid == pageurl {
		fmt.Println("❌ PAGEID was not found in NOTION_PAGEURL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
	}

	// CHECK PAGE CONTENT
	client := notionapi.NewClient(notionapi.Token(token))

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		fmt.Println("Failed retrieving page children blocks:", err)
		os.Exit(92)
	}

	if active, err := notionion.GetProxyStatus(children); err != nil {
		fmt.Println(err)
	} else if active {
		fmt.Println("📶 Proxy is active")
	} else {
		fmt.Println("📴 Proxy is inactive. Activate it by checking the \"OFF\" box")
	}

	// forward, err := notionion.RequestForwardButtonStatus(client, pageid)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if forward {
	// 	fmt.Println("📨 Forward request")
	// }

	// Request section checks
	if _, err := notionion.GetRequestBlock(children); err != nil {
		fmt.Println("❌ Request block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	} else {
		fmt.Println("➡️ Request block found")
	}
	if err := notionion.DisableRequestButtons(client, pageid); err != nil {
		fmt.Println(err)
	}

	codeReq, err := notionion.GetRequestCodeBlock(children)
	if err != nil {
		fmt.Println("❌ Request code block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	}
	notionion.ClearRequestCode(client, codeReq.ID)

	// Response section checks
	if _, err := notionion.GetResponseBlock(children); err != nil {
		fmt.Println("❌ Response block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	} else {
		fmt.Println("⬅️ Response block found")
	}

	codeResp, err := notionion.GetResponseCodeBlock(children)
	if err != nil {
		fmt.Println("❌ Response code block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	}
	notionion.ClearResponseCode(client, codeResp.ID)

	//PROXY SECTION
	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true

	// Request HTTP Handler
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if active, err := notionion.RequestProxyStatus(client, pageid); err != nil {
				fmt.Println(err)
				return r, nil
			} else if active {
				//reset response section
				notionion.ClearResponseCode(client, codeResp.ID)
				// Print request on Notion proxy page
				requestDump, err := httputil.DumpRequest(r, true)
				if err != nil {
					fmt.Println(err)
				}
				notionion.UpdateCodeContent(client, codeReq.ID, string(requestDump))
				//+enable button
				if err := notionion.EnableRequestButtons(client, pageid); err != nil {
					fmt.Println(err)
				}
				//wait for action (forward or drop)
				action := notionion.WaitAction(client, pageid)

				//disable button
				if err := notionion.DisableRequestButtons(client, pageid); err != nil {
					fmt.Println(err)
				}

				switch action {
				case notionion.FORWARD:
					//todo: retrieve code content -> to string
					reqFromPage, err := notionion.RequestRequestCodeContent(client, pageid)
					if err != nil {
						fmt.Println("Failed to retrieve request from notion proxy page:", err)
					}
					reader := bufio.NewReader(strings.NewReader(reqFromPage))
					if r, err = http.ReadRequest(reader); err != nil {
						fmt.Println("Failed parsing request from notion proxy page:", err)
					}
					return r, nil
				case notionion.DROP:
					return nil, nil
				}
			}
			return r, nil
		})

	// Response Handler
	proxy.OnResponse().DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if active, err := notionion.RequestProxyStatus(client, pageid); err != nil {
				fmt.Println(err)
				return resp
			} else if active {
				// Print response on Notion proxy page
				responseDump, err := httputil.DumpResponse(resp, true)
				if err != nil {
					fmt.Println(err)
				}

				// Print response in Notion proxy page
				notionion.UpdateCodeContent(client, codeResp.ID, string(responseDump))
			}
			return resp
		})

	fmt.Printf("🧅 Launch notionion proxy !\n\n")
	log.Fatal(http.ListenAndServe(":8080", proxy))

}
