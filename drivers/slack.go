package drivers

import (
	"log"
	"os"

	"github.com/akash/booker/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func ConnectToSlackViaSocketmode() (*socketmode.Client, error) {

	env := config.GetEnv()
	botToken := env["SLACK_BOT_TOKEN"]
	appToken := env["SLACK_APP_TOKEN"]

	api := slack.New(
		botToken,
		// slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	client := socketmode.New(
		api,
		// socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return client, nil
}
