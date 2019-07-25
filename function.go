package gcf_slack_sample

import (
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	threadTs slack.RTMsgOption
)

type SlackParams struct {
	accessToken string
	botUserID   string
	rtm         *slack.RTM
}

func PostMessage(w http.ResponseWriter, r *http.Request) { // 引数は飾り。Cloud Functionのhttp trigger で必要
	params := SlackParams{
		accessToken: os.Getenv("ACCESS_TOKEN"),
		botUserID:   "",
	}

	api := slack.New(params.accessToken)
	params.rtm = api.NewRTM()

	go params.rtm.ManageConnection()

	go func() {
		for msg := range params.rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				params.botUserID = ev.Info.User.ID
			case *slack.MessageEvent:
				if !strings.Contains(ev.Msg.Text, params.botUserID) {
					continue
				}
				threadTs = slack.RTMsgOptionTS(ev.ThreadTimestamp)
				params.rtm.SendMessage(params.rtm.NewOutgoingMessage("Sample TEST", ev.Channel, threadTs))
			}
		}
	}()
}
