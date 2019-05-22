package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

type Webhook struct {
	Instance       string `json:"instance"`
	CurrentApprove string `json:"current_approve"`
	Action         string `json:"action"`
	FromUser       string `json:"from_user"`
	ToUsers        string `json:"to_users"`
}

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var webhook Webhook
	if err := json.Unmarshal(buf.Bytes(), &webhook); err != nil {
		mlog.Error(err.Error())
		return
	}
	fmt.Print("webhook.action: ", webhook.Action)

	username := "sunhaolin"

	userID := ""
	if user, err := p.API.GetUserByUsername(username); err != nil {
		mlog.Error(err.Error())
		return
	} else {
		userID = user.Id
	}
	fmt.Print("webhook.userID: ", userID)
	fmt.Print("webhook.BotUserID: ", p.BotUserID)

	channel, err := p.API.GetDirectChannel(userID, p.BotUserID)
	if err != nil {
		mlog.Error("Couldn't get bot's DM channel", mlog.String("user_id", userID))
		return
	}

	post := &model.Post{
		UserId: p.BotUserID,
		Type:   "custom_webhook",
		Props: map[string]interface{}{
			"from_webhook": "true",
		},
	}

	message := fmt.Sprintf("New comment by %s", webhook.Action)

	channelID := channel.Id

	post.Message = message
	post.ChannelId = channelID

	fmt.Print("channelID: ", channelID)

	if _, err := p.API.CreatePost(post); err != nil {
		mlog.Error(err.Error())
	}
}
