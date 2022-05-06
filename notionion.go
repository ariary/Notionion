package main

import (
	"fmt"
	"log"
	"net/http"
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

	// Check page content
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

	forward, err := notionion.RequestForwardButtonStatus(client, pageid)
	if err != nil {
		fmt.Println(err)
	}
	if forward {
		fmt.Println("📨 Forward request")
	}

	// Request section checks
	if _, err := notionion.GetRequestBlock(children); err != nil {
		fmt.Println("❌ Request block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	} else {
		fmt.Println("➡️ Request block found")
	}

	codeReq, err := notionion.GetRequestCodeBlock(children)
	if err != nil {
		fmt.Println("❌ Request code block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	}

	// Response section checks

	if _, err := notionion.GetResponseBlock(children); err != nil {
		fmt.Println("❌ Response block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	} else {
		fmt.Println("⬅️ Response block found")
	}

	codeRes, err := notionion.GetResponseCodeBlock(children)
	if err != nil {
		fmt.Println("❌ Response code block not found in the proxy page")
		fmt.Println(err)
		os.Exit(92)
	}

	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true

	// Request Handler
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			fmt.Println("Youhou")
			if active, err := notionion.RequestProxyStatus(client, pageid); err != nil {
				fmt.Println(err)
				return r, nil
			} else if active {
				// Print request on Notion proxy page
				notionion.UpdateCodeContent(client, codeReq.ID, r.Host) //todo: request to string and string -> request
				//wait for action (forward or drop)
				action := notionion.WaitAction()
				notionion.UpdateCodeContent(client, codeRes.ID, r.Host) //todo: request to string and string -> request
				switch action {
				case notionion.FORWARD:
					//todo: retrieve code content -> to string
					return r, nil
				case notionion.DROP:
					return nil, nil
				}
			}
			return r, nil
		})
	fmt.Println("🧅 Launch notionion proxy !")
	log.Fatal(http.ListenAndServe(":8080", proxy))
	// //CODEBLOCK UPDATE
	// codeRes, err := notionion.GetResponseCodeBlock(children)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// _, err = notionion.UpdateCodeContent(client, codeRes.ID, "this is a test")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//BUTTON DISABLING
	// _, err = notionion.GetRequestButtonsColumnBlock(children)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if err := notionion.DisableRequestButtons(client, pageid); err != nil {
	// 	fmt.Println(err)
	// }
	// button, _ := notionion.RequestRequestButtonByName(client, pageid, notionion.FORWARD)
	// fmt.Printf("%+v", button.ToDo.RichText[0].Annotations)

}
