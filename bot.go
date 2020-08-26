package main

import (
	"github.com/Krognol/go-wolfram"
	"github.com/christianrondeau/go-wit"
	"fmt"
	 "strings"
	 "os"
	 "github.com/slack-go/slack"
)
var slackClient *slack.Client
var witClient *wit.Client
var rtm *slack.RTM
var wolframClient *wolfram.Client
const confidenceThreshold = 0.5
func main() {
	fmt.Print("Hi!!")
	var slackToken string = os.Getenv("SLACK_TOKEN")
	var witToken string = os.Getenv("WIT_TOKEN")
	slackClient = slack.New(slackToken)
	rtm = slackClient.NewRTM()
	witClient = wit.NewClient(witToken)
	wolframClient = &wolfram.Client{AppID: os.Getenv("WOLFRAM_TOKEN")}
	
	go rtm.ManageConnection()

	for{ 
		select {
		case message := <-rtm.IncomingEvents:
			fmt.Print("Event type: ")
			//event := message.Data.(type)
			switch event := message.Data.(type){
			case *slack.ConnectedEvent:
				fmt.Println("New Connection, current connection count: ", event.ConnectionCount)
			case *slack.MessageEvent:
				
				//fmt.Println("Message receieved, message: ", event)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if event.User != info.User.ID && strings.HasPrefix(event.Text, prefix) {
					go handleNewMessage(event, event.Channel);
					//rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", event.Channel))
				}
			case *slack.RTMError:
				fmt.Println("Error: ", event.Error())
			case *slack.InvalidAuthEvent:
				fmt.Println("Invalid credentials")
			default:
				fmt.Println("Unexpected")
			}
			// if event == slack.ConnectedEvent { 
			// 	fmt.Println("New Connection, current connection count: ", event.ConnectionCount)
			// } else if event == slack.MessageEvent { 
			// 	fmt.Println("Message receieved, message: ", event)
			// } else if event == slack.RTMError { 
			// 	fmt.Println("Error: ", event.Error())
			// } else if event == slack.InvalidAuthEvent { 
			// 	fmt.Println("Invalid credentials")
			// }
		}
	}
}

func handleNewMessage(ev *slack.MessageEvent, channel string) {
	res, err := witClient.Message(ev.Msg.Text)
	if err != nil {
		fmt.Printf("wit.ai error: %v", err)
		return
	}

	var topEntity wit.MessageEntity
	var	topEntityKey string
	

	for key, entities := range res.Entities {
		for _, entity := range entities {
			if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
				topEntity = entity
				topEntityKey = key
			}
		}
	}

	sendResponse(ev, topEntity, topEntityKey, channel)
}
func sendResponse(ev *slack.MessageEvent, topEntity wit.MessageEntity, topEntityKey string, channel string) { 
	switch topEntityKey { 
	case "greetings":
		rtm.SendMessage(rtm.NewOutgoingMessage("Hey there! How can I help you?", channel))
	case "wolfram_search_query":
		res, err := wolframClient.GetShortAnswerQuery(topEntity.Value.(string), wolfram.Metric, 1000)
		if err != nil { 
			fmt.Printf("Unable to get wolfram result: %v", err)
			return
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(res, channel))
	default:
		rtm.SendMessage(rtm.NewOutgoingMessage("Invalid Message! Please try again.", channel))

	}
}
