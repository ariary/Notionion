package notionion

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jomei/notionapi"
)

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
