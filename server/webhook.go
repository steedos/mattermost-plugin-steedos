package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

// WorkflowWebhook apps发送的webhook结构
type WorkflowWebhook struct {
	Instance       Instance `json:"instance"`
	CurrentApprove Approve  `json:"current_approve"`
	Action         string   `json:"action"`
	FromUser       User     `json:"from_user"`
	ToUsers        []User   `json:"to_users"`
}

// Instance 申请单
type Instance struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Space       string `json:"space"`
	RedirectURL string `json:"redirectUrl"`
}

// Approve 历史步骤
type Approve struct {
	ID string `json:"_id"`
}

// User 用户结构
type User struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var webhook WorkflowWebhook
	if err := json.Unmarshal(buf.Bytes(), &webhook); err != nil {
		mlog.Error(err.Error())
		return
	}
	fmt.Print("webhook.action: ", webhook.Action)
	fmt.Print("webhook: ", webhook)

	for _, u := range webhook.ToUsers {
		fmt.Print("toUsers.Username: ", u.Username)

		userID := ""
		if user, err := p.API.GetUserByUsername(u.Username); err != nil {
			mlog.Error(err.Error())
			continue
		} else {
			userID = user.Id
		}

		channel, err := p.API.GetDirectChannel(userID, p.BotUserID)
		if err != nil {
			mlog.Error("Couldn't get bot's DM channel", mlog.String("user_id", userID))
			continue
		}

		post := &model.Post{
			UserId: p.BotUserID,
			Type:   "custom_webhook",
			Props: map[string]interface{}{
				"from_webhook": "true",
			},
		}

		fmt.Print("webhook.http.Request-URL: ", r.URL)
		fmt.Print("webhook.http.Request-RequestURI: ", r.RequestURI)
		fmt.Print("webhook.http.Request-Host: ", r.Host)
		fmt.Print("webhook.http.Request-Proto: ", r.Proto)

		message := fmt.Sprintf("请确认: [%s](%s)", webhook.Instance.Name, webhook.Instance.RedirectURL)
		if "engine_submit" == webhook.Action {
			message = fmt.Sprintf("请审批: [%s](%s)", webhook.Instance.Name, webhook.Instance.RedirectURL)
		}

		post.Message = message
		post.ChannelId = channel.Id

		fmt.Print("channelID: ", channel.Id)

		if _, err := p.API.CreatePost(post); err != nil {
			mlog.Error(err.Error())
		}
	}

}
