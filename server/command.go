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
* |/steedos subscribe list| - Will list the current channel subscriptions
* |/steedos subscribe objectName| - Subscribe the current channel to receive notifications about create update and delete for an object
* |/steedos unsubscribe objectName| - Unsubscribe the current channel from an object
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
	parameters := []string{}
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if command != "/steedos" {
		return &model.CommandResponse{}, nil
	}

	switch action {
	case "subscribe":
		objectName := ""
		txt := ""
		if len(parameters) == 0 {
			return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please specify an objectName or 'list' command."), nil
		} else if len(parameters) == 1 && parameters[0] == "list" {
			subs, err := p.GetSubscriptionsByChannel(args.ChannelId)
			if err != nil {
				return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), nil
			}

			if len(subs) == 0 {
				txt = "Currently there are no subscriptions in this channel"
			} else {
				txt = "### Subscriptions in this channel\n"
			}
			for _, sub := range subs {
				txt += fmt.Sprintf("* `%s`\n", sub.ObjectName)
			}
			return p.getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, txt), nil
		}

		objectName = parameters[0]

		if err := p.Subscribe(context.Background(), args.UserId, args.ChannelId, objectName); err != nil {
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
