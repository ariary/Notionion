package notionion

import (
	"context"
	"fmt"
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
func GetProxyStatus(children notionapi.Blocks) (bool, error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "to_do" {
			todo := children[i].(*notionapi.ToDoBlock).ToDo
			if todo.RichText[0].Text.Content == "ON" {
				return todo.Checked, nil
			}
		}
	}
	err := fmt.Errorf("Failed retrieving proxy status button")
	return false, err
}

//GetRequestBlock: retrieve "Request" block from page's blocks
func GetRequestBlock(children notionapi.Blocks) (requestBlock notionapi.Heading2Block, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "heading_2" {
			headingBlock := children[i].(*notionapi.Heading2Block)
			if strings.Contains(headingBlock.Heading2.RichText[0].Text.Content, "Request") {
				return *headingBlock, nil
			}
		}
	}
	err = fmt.Errorf("Failed retrieving \"request\" section")
	return requestBlock, err
}

//GetResponseBlock: retrieve "Response" block from page's blocks
func GetResponseBlock(children notionapi.Blocks) (responseBlock notionapi.Heading2Block, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "heading_2" {
			headingBlock := children[i].(*notionapi.Heading2Block)
			if strings.Contains(headingBlock.Heading2.RichText[0].Text.Content, "Response") {
				return *headingBlock, nil
			}
		}
	}
	err = fmt.Errorf("Failed retrieving \"response\" section")
	return responseBlock, err
}

//GetCodeBlockByName: Obtain the code block object under the section specified by name (name={"Request","Response"})
func GetCodeBlockByName(children notionapi.Blocks, name string) (requestCodeBlock notionapi.CodeBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "code" {
			codeBlock := children[i].(*notionapi.CodeBlock)
			if i > 0 {
				if children[i-1].GetType() == "heading_2" {
					above := children[i-1].(*notionapi.Heading2Block)
					if strings.Contains(above.Heading2.RichText[0].Text.Content, name) {
						return *codeBlock, nil
					}
				}
			}
		}
	}
	err = fmt.Errorf("Failed retrieving request code block within \"request\" section")
	return requestCodeBlock, err
}

//GetRequestCodeBlock: Obtain the code block object under the request heading
func GetRequestCodeBlock(children notionapi.Blocks) (requestCodeBlock notionapi.CodeBlock, err error) {
	return GetCodeBlockByName(children, "Request")
}

//GetResponseCodeBlock: Obtain the code block object under the response heading
func GetResponseCodeBlock(children notionapi.Blocks) (requestCodeBlock notionapi.CodeBlock, err error) {
	return GetCodeBlockByName(children, "Response")
}

//UpdateCodeContent: update code block with content
func UpdateCodeContent(client *notionapi.Client, codeBlockID notionapi.BlockID, content string) (notionapi.Block, error) {
	//construct code block containing request
	code := notionapi.CodeBlock{
		Code: notionapi.Code{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: content,
					},
					Annotations: &notionapi.Annotations{
						Bold:          false,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          false,
						Color:         "",
					},
				},
			},
			Language: "html",
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Code: &code.Code,
	}

	return client.Block.Update(context.Background(), codeBlockID, updateReq)
}

//GetRequestButtonsColumnBlock: retrieve buttons within request block (column list block)
func GetRequestButtonsColumnBlock(children notionapi.Blocks) (buttonsBlock notionapi.ColumnListBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == "column_list" {
			buttonsBlock := children[i].(*notionapi.ColumnListBlock)
			if buttonsBlock.HasChildren && i > 0 && children[i-1].GetType() == "code" { //TODO: check if the code is under request heading
				return *buttonsBlock, nil
			}
		}
	}
	err = fmt.Errorf("GetRequestButtonsColumnBlock: failed to retrieve column block containg button within request section")
	return buttonsBlock, err
}

//RequestRequestButtonByName:return specific to_do block within "request" block is checked.
// name: {"FORWARD", "DROP"}
func RequestRequestButtonByName(client *notionapi.Client, pageid string, name string) (button notionapi.ToDo, err error) {
	children, err := RequestProxyPageChildren(client, pageid)
	if err != nil {
		return button, err
	}
	buttonsBlock, err := GetRequestButtonsColumnBlock(children)
	if err != nil {
		return button, err
	}

	columnsList, err := client.Block.GetChildren(context.Background(), buttonsBlock.ID, nil)
	if err != nil {
		return button, err
	}
	columns := columnsList.Results
	for i := 0; i < len(columns); i++ {
		buttonsList, err := client.Block.GetChildren(context.Background(), columns[i].GetID(), nil)
		if err != nil {
			return button, err
		}

		for j := 0; j < len(buttonsList.Results); j++ {
			if buttonsList.Results[j].GetType() == "to_do" {
				todo := buttonsList.Results[j].(*notionapi.ToDoBlock).ToDo
				if todo.RichText[0].Text.Content == name {
					return todo, err
				}
			}
		}
	}

	return button, err
}

//RequestForwardButtonStatus: check if forward button is checked
func RequestForwardButtonStatus(client *notionapi.Client, pageid string) (checked bool, err error) {
	forward, err := RequestRequestButtonByName(client, pageid, "FORWARD")
	if err != nil {
		return false, err
	}

	return forward.Checked, err
}

//RequestDropButtonStatus: check if drop button is checked
func RequestDropButtonStatus(client *notionapi.Client, pageid string) (checked bool, err error) {
	drop, err := RequestRequestButtonByName(client, pageid, "DROP")
	if err != nil {
		return false, err
	}

	return drop.Checked, err
}

func DisableRequestButtons(client *notionapi.Client, pageid string) error {
	// forward, err := RequestRequestButtonByName(client, pageid, "FORWARD")
	// if err != nil {
	// 	return err
	// }

	//make function requestbutton return button
	//getstatus use this function and chech "checked" status
	return nil
}

// func GetRequestParagraphBlock(children notionapi.Blocks) (requestParagraphBlock notionapi.ParagraphBlock) {
// 	for i := 0; i < len(children); i++ {
// 		if children[i].GetType() == "paragraph" {
// 			paragraphBlock := children[i].(*notionapi.ParagraphBlock)
// 			if i > 0 {
// 				if children[i-1].GetType() == "heading_2" {
// 					above := children[i-1].(*notionapi.Heading2Block)
// 					if strings.Contains(above.Heading2.RichText[0].Text.Content, "Request") {
// 						return *paragraphBlock
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return requestParagraphBlock
// }

//UpdateRequestContent: update text in paragraph block within request section
// func UpdateRequestContent(client *notionapi.Client, requestCodeBlockID notionapi.BlockID, request string) (notionapi.Block, error) {
// 	//construct paragraph block containing request
// 	paragraph := notionapi.ParagraphBlock{
// 		Paragraph: notionapi.Paragraph{
// 			RichText: []notionapi.RichText{
// 				{
// 					Type: notionapi.ObjectTypeText,
// 					Text: notionapi.Text{
// 						Content: request,
// 					},
// 					Annotations: &notionapi.Annotations{
// 						Bold:          true,
// 						Italic:        false,
// 						Strikethrough: false,
// 						Underline:     false,
// 						Code:          true,
// 						Color:         "",
// 					},
// 				},
// 			},
// 		},
// 	}
// 	//AppendBlockChildrenRequest
// 	updateReq := &notionapi.BlockUpdateRequest{
// 		Paragraph: &paragraph.Paragraph,
// 	}

// 	// send update request
// 	return client.Block.Update(context.Background(), requestCodeBlockID, updateReq)
// }
