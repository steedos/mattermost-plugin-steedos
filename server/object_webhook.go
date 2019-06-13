package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

// ObjectWebhook creator发送的webhook结构
type ObjectWebhook struct {
	Data              map[string]interface{} `json:"data"`
	Action            string                 `json:"action"`
	ActionUserInfo    ActionUserInfo         `json:"actionUserInfo"`
	ObjectName        string                 `json:"objectName"`
	ObjectDisplayName string                 `json:"objectDisplayName"`
	NameFieldKey      string                 `json:"nameFieldKey"`
	RedirectURL       string                 `json:"redirectUrl"`
}

// ActionUserInfo 用户结构
type ActionUserInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

func (p *Plugin) handleObjectWebhook(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var webhook ObjectWebhook
	if err := json.Unmarshal(buf.Bytes(), &webhook); err != nil {
		mlog.Error(err.Error())
		return
	}
	fmt.Print("webhook.action》》》》》》》》》》: ", webhook.Action)
	fmt.Print("webhook》》》》》》》》》》: ", webhook)

	subs := p.GetSubscribedChannelsForObjectName(webhook.ObjectName)

	if subs == nil || len(subs) == 0 {
		fmt.Print("111111111111111111111111: ")
		return
	}

	message := ""
	if "create" == webhook.Action {
		message = fmt.Sprintf("%s新建了%s: [%s](%s)", webhook.ActionUserInfo.Name, webhook.ObjectDisplayName, webhook.Data[webhook.NameFieldKey], webhook.RedirectURL)
	} else if "update" == webhook.Action {
		message = fmt.Sprintf("%s更新了%s: [%s](%s)", webhook.ActionUserInfo.Name, webhook.ObjectDisplayName, webhook.Data[webhook.NameFieldKey], webhook.RedirectURL)
	} else if "delete" == webhook.Action {
		message = fmt.Sprintf("%s删除了%s: [%s](%s)", webhook.ActionUserInfo.Name, webhook.ObjectDisplayName, webhook.Data[webhook.NameFieldKey], webhook.RedirectURL)
	}
	fmt.Print("222222222222222222222222222: ")

	actionUserInfo, _ := json.Marshal(webhook.ActionUserInfo)
	var actionUserInfoMapResult map[string]interface{}
	json.Unmarshal(actionUserInfo, &actionUserInfoMapResult)

	fmt.Print("_id_id_id_id_id_id_id_id_id: ", webhook.Data["_id"])

	dataMap := map[string]interface{}{
		"_id": webhook.Data["_id"],
	}

	dataMap[webhook.NameFieldKey] = webhook.Data[webhook.NameFieldKey]

	post := &model.Post{
		UserId: p.BotUserID,
		Type:   "custom_object_webhook",
		Props: map[string]interface{}{
			"from_object_webhook": "true",
			"data":                dataMap,
			"action":              webhook.Action,
			"actionUserInfo":      actionUserInfoMapResult,
			"objectName":          webhook.ObjectName,
			"objectDisplayName":   webhook.ObjectDisplayName,
			"nameFieldKey":        webhook.NameFieldKey,
			"redirectUrl":         webhook.RedirectURL,
		},
	}

	post.Message = message

	for _, sub := range subs {
		fmt.Print("4444444444444444444444444444: ", sub.ChannelID)

		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			mlog.Error(err.Error())
		}
	}

}
