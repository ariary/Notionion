package notionion

import (
	"context"
	"strings"

	"github.com/jomei/notionapi"
)

//RequestProxyPageChildren: Returns the children block of the Listener page
func RequestProxyPageChildren(client *notionapi.Client, pageid string) (childrenBlocks notionapi.Blocks, err error) {
	children, err := client.Block.GetChildren(context.Background(), notionapi.BlockID(pageid), nil)
	return children.Results, err
}

//RequestProxyStatus: request notion api to determine if proxy is active
func RequestProxyStatus(client *notionapi.Client, pageid string) (active bool, err error) {
	children, err := RequestProxyPageChildren(client, pageid)
	if err != nil {
		return false, err
	}
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "to_do" {
			todo := children[i].(*notionapi.ToDoBlock).ToDo
			if todo.RichText[0].Text.Content == "OFF" {
				return todo.Checked, err
			}
		}
	}
	return false, nil
}

//GetProxyStatus: get proxy status from page's blocks
func GetProxyStatus(children notionapi.Blocks) bool {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "to_do" {
			todo := children[i].(*notionapi.ToDoBlock).ToDo
			if todo.RichText[0].Text.Content == "ON" {
				return todo.Checked
			}
		}
	}
	return false
}

//GetRequestBlock: retrieve "Request" block from page's blocks
func GetRequestBlock(children notionapi.Blocks) (requestBlock notionapi.Heading2Block) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "heading_2" {
			headingBlock := children[i].(*notionapi.Heading2Block)
			if strings.Contains(headingBlock.Heading2.RichText[0].Text.Content, "Request") {
				return *headingBlock
			}
		}
	}
	return requestBlock
}

//GetResponseBlock: retrieve "Response" block from page's blocks
func GetResponseBlock(children notionapi.Blocks) (requestBlock notionapi.Heading2Block) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "heading_2" {
			headingBlock := children[i].(*notionapi.Heading2Block)
			if strings.Contains(headingBlock.Heading2.RichText[0].Text.Content, "Response") {
				return *headingBlock
			}
		}
	}
	return requestBlock
}

func GetRequestParagraphBlock(children notionapi.Blocks) (requestParagraphBlock notionapi.ParagraphBlock) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "paragraph" {
			paragraphBlock := children[i].(*notionapi.ParagraphBlock)
			if i > 0 {
				if children[i-1].GetType() == "heading_2" {
					above := children[i-1].(*notionapi.Heading2Block)
					if strings.Contains(above.Heading2.RichText[0].Text.Content, "Request") {
						return *paragraphBlock
					}
				}
			}
		}
	}
	return requestParagraphBlock
}

//UpdateRequestContent: update text in paragraph block within request section
func UpdateRequestContent(client *notionapi.Client, requestCodeBlockID notionapi.BlockID, request string) (notionapi.Block, error) {
	//construct paragraph block containing request
	paragraph := notionapi.ParagraphBlock{
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: request,
					},
					Annotations: &notionapi.Annotations{
						Bold:          true,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          true,
						Color:         "",
					},
				},
			},
		},
	}
	//AppendBlockChildrenRequest
	updateReq := &notionapi.BlockUpdateRequest{
		Paragraph: &paragraph.Paragraph,
	}

	// send update request
	return client.Block.Update(context.Background(), requestCodeBlockID, updateReq)
}

//UpdateResponseContent: update response block

//CheckRequestSendingBox: Check if the to_do block to send request is checked


