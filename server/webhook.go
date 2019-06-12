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
	ID            string `json:"_id"`
	Name          string `json:"name"`
	Space         string `json:"space"`
	RedirectURL   string `json:"redirectUrl"`
	ApplicantName string `json:"applicant_name"`
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
			Type:   "custom_workflow_webhook",
			Props: map[string]interface{}{
				"from_webhook":   "true",
				"action":         webhook.Action,
				"name":           webhook.Instance.Name,
				"redirectUrl":    webhook.Instance.RedirectURL,
				"applicant_name": webhook.Instance.ApplicantName,
			},
		}

		message := fmt.Sprintf("你有新的待办文件: [%s](%s), 提交人: %s", webhook.Instance.Name, webhook.Instance.RedirectURL, webhook.Instance.ApplicantName)

		post.Message = message
		post.ChannelId = channel.Id

		fmt.Print("channelID: ", channel.Id)

		if _, err := p.API.CreatePost(post); err != nil {
			mlog.Error(err.Error())
		}
	}

}
