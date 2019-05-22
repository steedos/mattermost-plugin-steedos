package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/plugin"

	"github.com/mattermost/mattermost-server/model"
)

const (
	STEEDOS_USERNAME = "Steedos Plugin"
)

const COMMAND_HELP = `
* |/steedos subscribe| - Subscribe the current channel to receive notifications
* |/steedos unsubscribe| - Unsubscribe the current channel
`

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "steedos",
		DisplayName:      "Steedos",
		Description:      "Integration with Steedos.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: subscribe, unsubscribe, help",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     STEEDOS_USERNAME,
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = split[1]
	}

	if command != "/steedos" {
		return &model.CommandResponse{}, nil
	}

	switch action {
	case "subscribe":
		if err := p.Subscribe(context.Background(), args.UserId, args.ChannelId); err != nil {
			return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), nil
		}

		return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Successfully subscribed.")), nil
	case "unsubscribe":

		if err := p.Unsubscribe(args.ChannelId); err != nil {
			mlog.Error(err.Error())
			return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Encountered an error trying to unsubscribe. Please try again."), nil
		}

		return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Succesfully unsubscribed.")), nil

	case "help":
		text := "###### Mattermost Steedos Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	case "":
		text := "###### Mattermost Steedos Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil

	}

	return &model.CommandResponse{}, nil
}
