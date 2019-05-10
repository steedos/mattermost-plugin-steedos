package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	steedosClient *Client

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

type UserInfo struct {
	UserID    string `json:"userId"`
	AuthToken string `json:"authToken"`
}

type Body struct {
}

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		return err
	}

	p.steedosClient = NewClient(config.APIURL, config.APIKey, config.APISecret)

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	switch path := r.URL.Path; path {
	case "/startup":
		p.handleStartup(c, w, r)
	default:
		http.NotFound(w, r)
	}

}

func (p *Plugin) handleStartup(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Print(">>>>>>>>>>c.SessionId： ", c.SessionId)

	config := p.getConfiguration()

	userId := r.Header.Get("Mattermost-User-Id")
	if userId == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var user *model.User
	var err *model.AppError
	user, err = p.API.GetUser(userId)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
	}

	fmt.Print(">>>>>>>user.Username: ", user.Username)
	body := Body{}

	var ret UserInfo
	erro := p.steedosClient.request("GET", fmt.Sprintf("/api/jwt/sso"), body, &ret, user.Username, c.SessionId)
	if erro != nil {
		http.Error(w, erro.Error(), erro.StatusCode)
		return
	}

	jsonData, _ := json.Marshal(ret)
	var mapResult map[string]interface{}
	json.Unmarshal(jsonData, &mapResult)
	mapResult["url"] = config.APIURL
	jsonStr, _ := json.Marshal(mapResult)
	fmt.Print(">>>>>>>>>>jsonData： ", string(jsonStr))

	fmt.Fprint(w, string(jsonStr))

}

// See https://developers.mattermost.com/extend/plugins/server/reference/
