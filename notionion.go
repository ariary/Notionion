package main

import (
	"flag"
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
	port := "8080"
	flag.Parse()
	if len(flag.Args()) > 0 {
		port = flag.Arg(0)
	}
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
	proxy.OnRequest().Do(notionion.ProxyRequestHTTPHandler(client, pageid, codeReq, codeResp))

	// Response Handler
	proxy.OnResponse().Do(notionion.ProxyResponseHTTPHandler(client, pageid, codeResp))

	fmt.Printf("🧅 Launch notionion proxy on port %s !\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))

}
