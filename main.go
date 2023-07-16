package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akash/booker/config"
	"github.com/akash/booker/drivers"
	"github.com/akash/booker/handlers"
	"github.com/akash/booker/model"
	"github.com/akash/booker/views"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func init() {
	config.LoadEnv()
	model.InitDB()
}

func main() {

	client, err := drivers.ConnectToSlackViaSocketmode()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("DEBUG: Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("DEBUG: Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("DEBUG: Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				fmt.Printf("DEBUG: Event received =>  %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
						if err != nil {
							fmt.Printf("failed posting message: %v", err)
						}
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
					case *slackevents.AppHomeOpenedEvent:
						handlers.AppHomeOpenedEvent(ev, client)
					}
				default:
					client.Debugf("DEBUG: unsupported Events API event received")
				}

			case socketmode.EventTypeInteractive:

				callback, ok := evt.Data.(slack.InteractionCallback)
				blockActions := callback.ActionCallback.BlockActions

				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				var payload interface{}

				switch callback.Type {

				case slack.InteractionTypeBlockActions:
					if blockActions != nil {

						fmt.Println("Action ID: ", blockActions[0].ActionID,
							"\n Action.Value:", blockActions[0].Value)

						switch blockActions[0].ActionID {
						case views.ScheduleBookingActionID:
							fmt.Println("DEBUG: block Actions Triggered", views.ScheduleBookingActionID,
								callback.TriggerID, client)
							handlers.AppHomeScheduledBookingModal(callback.TriggerID, client)

						case views.CurrentStatusActionID:
							fmt.Println("DEBUG: block Actions Triggered", views.CurrentStatusActionID)
							handlers.PublishCurrentRoomStatus(callback, client)

						case views.AppHomeCancelBookingActionID:
							fmt.Println("DEBUG: block Actions Triggered", views.AppHomeCancelBookingActionID)
							handlers.AppHomeDeclineBooking(blockActions[0].Value, callback, client)

						default:
							fmt.Println("undeclared block action: ", blockActions[0])
						}

					} else {
						fmt.Println("Block Actions Empty")
					}

				case slack.InteractionTypeViewSubmission:
					fmt.Println("DEBUG:", "InteractionTypeViewSubmission")
					switch callback.View.CallbackID {
					case views.ScheduleBookingModalCallbackID:
						handlers.PublishScheduleBooking(callback, client)

					}
					// handlers.PublishStickie(callback, client)

				case slack.InteractionTypeShortcut:
				// case slack.InteractionTypeViewSubmission:
				// 	// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:
					client.Debugf("DEBUG: unsupported Interactive Events API event received")
				}

				client.Ack(*evt.Request, payload)
			// case socketmode.EventTypeSlashCommand:
			// 	cmd, ok := evt.Data.(slack.SlashCommand)
			// 	if !ok {
			// 		fmt.Printf("Ignored %+v\n", evt)

			// 		continue
			// 	}

			// 	client.Debugf("Slash command received: %+v", cmd)

			// 	payload := map[string]interface{}{
			// 		"blocks": []slack.Block{
			// 			slack.NewSectionBlock(
			// 				&slack.TextBlockObject{
			// 					Type: slack.MarkdownType,
			// 					Text: "foo",
			// 				},
			// 				nil,
			// 				slack.NewAccessory(
			// 					slack.NewButtonBlockElement(
			// 						"",
			// 						"somevalue",
			// 						&slack.TextBlockObject{
			// 							Type: slack.PlainTextType,
			// 							Text: "bar",
			// 						},
			// 					),
			// 				),
			// 			),
			// 		},
			// 	}

				client.Ack(*evt.Request, payload)
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}
