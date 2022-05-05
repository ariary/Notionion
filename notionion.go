package main

import (
	"fmt"
	"os"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

func main() {
	// integration token
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("❌ Please set NOTION_TOKEN envvar with your integration token before launching notionion")
		os.Exit(92)
	}
	pageid := os.Getenv("NOTION_PAGEID")
	if pageid == "" {
		fmt.Println("❌ Please set NOTION_PAGEID envvar with your page id before launching notionion (CTRL+L on desktop app)")
		os.Exit(92)
	}

	client := notionapi.NewClient(notionapi.Token(token))
	// page, err := client.Page.Get(context.Background(), notionapi.PageID(pageid))
	// if err != nil {
	// 	fmt.Println("failed retrieving page:", pageid)
	// 	os.Exit(92)
	// }

	// children, err := notionion.GetProxyPageChildren(client, pageid)
	// if err != nil {
	// 	fmt.Println("Failed retrieving page children blocks:", err)
	// 	os.Exit(92)
	// }

	active, err := notionion.GetProxyStatus(client, pageid)
	if err != nil {
		fmt.Println("Failed retrieving proxy status")
		os.Exit(92)
	}

	if active {
		fmt.Println("📶 Proxy is active")
	} else {
		fmt.Println("📴 Proxy is inactive. Activate it by checking the \"OFF\" box")
	}

	// for i := 0; i < len(children); i++ {
	// 	c, err := client.Block.Get(context.Background(), children[i].GetID())
	// 	if err != nil {
	// 		fmt.Println("tttt")
	// 	}
	// 	fmt.Printf("%+v", c)
	// 	fmt.Println()
	// }

}
