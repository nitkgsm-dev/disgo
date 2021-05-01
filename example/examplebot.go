package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

const red = 16711680
const orange = 16562691
const green = 65280

const guildID = "817327181659111454"
const adminRoleID = "817327279583264788"
const testRoleID = "825156597935243304"

var logger = logrus.New()
var client = http.DefaultClient

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("starting TestBot...")

	dgo, err := disgo.NewBuilder(os.Getenv("token")).
		SetLogger(logger).
		SetRawGatewayEventsEnabled(true).
		SetHTTPClient(client).
		SetGatewayIntents(api.GatewayIntentsGuilds | api.GatewayIntentsGuildMessages | api.GatewayIntentsGuildMembers).
		SetMemberCachePolicy(api.MemberCachePolicyAll).
		AddEventListeners(&events.ListenerAdapter{
			OnRawGateway:         rawGatewayEventListener,
			OnGuildAvailable:     guildAvailListener,
			OnGuildMessageCreate: messageListener,
			OnSlashCommand:       slashCommandListener,
			OnButtonClick:        buttonClickListener,
		}).
		Build()
	if err != nil {
		logger.Fatalf("error while building disgo instance: %s", err)
		return
	}

	rawCmds := []*api.CommandCreate{
		{
			Name:              "eval",
			Description:       "runs some go code",
			DefaultPermission: ptrBool(false),
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeString,
					Name:        "code",
					Description: "the code to eval",
					Required:    true,
				},
			},
		},
		{
			Name:              "test",
			Description:       "test test test test test test",
			DefaultPermission: ptrBool(false),
		},
		{
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: ptrBool(false),
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeString,
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
			},
		},
		{
			Name:              "addrole",
			Description:       "This command adds a role to a member",
			DefaultPermission: ptrBool(false),
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to add a role to",
					Required:    true,
				},
				{
					Type:        api.CommandOptionTypeRole,
					Name:        "role",
					Description: "The role to add to a member",
					Required:    true,
				},
			},
		},
		{
			Name:              "removerole",
			Description:       "This command removes a role from a member",
			DefaultPermission: ptrBool(false),
			Options: []*api.CommandOption{
				{
					Type:        api.CommandOptionTypeUser,
					Name:        "member",
					Description: "The member to removes a role from",
					Required:    true,
				},
				{
					Type:        api.CommandOptionTypeRole,
					Name:        "role",
					Description: "The role to removes from a member",
					Required:    true,
				},
			},
		},
	}

	// using the api.RestClient directly to avoid the guild needing to be cached
	cmds, err := dgo.RestClient().SetGuildCommands(dgo.ApplicationID(), guildID, rawCmds...)
	if err != nil {
		logger.Errorf("error while registering guild commands: %s", err)
	}

	var cmdsPermissions []*api.SetGuildCommandPermissions
	for _, cmd := range cmds {
		var perms *api.CommandPermission
		if cmd.Name == "eval" {
			perms = &api.CommandPermission{
				ID:         adminRoleID,
				Type:       api.CommandPermissionTypeRole,
				Permission: true,
			}
		} else {
			perms = &api.CommandPermission{
				ID:         testRoleID,
				Type:       api.CommandPermissionTypeRole,
				Permission: true,
			}
		}
		cmdsPermissions = append(cmdsPermissions, &api.SetGuildCommandPermissions{
			ID:          cmd.ID,
			Permissions: []*api.CommandPermission{perms},
		})
	}
	if _, err = dgo.RestClient().SetGuildCommandsPermissions(dgo.ApplicationID(), guildID, cmdsPermissions...); err != nil {
		logger.Errorf("error while setting command permissions: %s", err)
	}

	err = dgo.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to discord: %s", err)
	}

	defer dgo.Close()

	logger.Infof("TestBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func guildAvailListener(event *events.GuildAvailableEvent) {
	logger.Printf("guild loaded: %s", event.Guild.ID)
}

func rawGatewayEventListener(event *events.RawGatewayEvent) {
	if event.Type == api.GatewayEventInteractionCreate {
		println(string(event.RawPayload))
	}
}

func buttonClickListener(event *events.ButtonClickEvent) {
	switch event.CustomID {
	case "test":
		err := event.Reply(&api.InteractionResponse{
			Type: api.InteractionResponseTypeButtonResponse,
			Data: &api.ButtonResponseData{
				Content: nil,
				//Components:      nil,
			},
		})
		if err != nil {
			logger.Errorf("error sending interaction response: %s", err)
		}
	}
}

func slashCommandListener(event *events.SlashCommandEvent) {
	switch event.CommandName {
	case "eval":
		go func() {
			code := event.Option("code").String()
			embed := api.NewEmbedBuilder().
				SetColor(orange).
				AddField("Status", "...", true).
				AddField("Time", "...", true).
				AddField("Code", "```go\n"+code+"\n```", false).
				AddField("Output", "```\n...\n```", false)
			_ = event.Reply(api.NewCommandResponseBuilder().SetEmbeds(embed.Build()).Build())

			start := time.Now()
			output, err := gval.Evaluate(code, map[string]interface{}{
				"disgo": event.Disgo(),
				"dgo":   event.Disgo(),
				"event": event,
			})

			elapsed := time.Since(start)
			embed.SetField(1, "Time", strconv.Itoa(int(elapsed.Milliseconds()))+"ms", true)

			if err != nil {
				_, _ = event.EditOriginal(api.NewFollowupMessageBuilder().
					SetEmbeds(embed.
						SetColor(red).
						SetField(0, "Status", "Failed", true).
						SetField(3, "Output", "```"+err.Error()+"```", false).
						Build(),
					).
					Build(),
				)
				return
			}
			_, err = event.EditOriginal(api.NewFollowupMessageBuilder().
				SetEmbeds(embed.
					SetColor(green).
					SetField(0, "Status", "Success", true).
					SetField(3, "Output", "```"+fmt.Sprintf("%+v", output)+"```", false).
					Build(),
				).
				Build(),
			); if err != nil {
				logger.Errorf("error sending interaction response: %s", err)
			}
		}()

	case "say":
		_ = event.Reply(api.NewCommandResponseBuilder().
			SetContent(event.Option("message").String()).
			SetAllowedMentionsEmpty().
			Build(),
		)

	case "test":
		if err := event.Reply(api.NewCommandResponseBuilder().
			SetEphemeral(true).
			SetContent("test1").
			SetEmbeds(api.NewEmbedBuilder().SetDescription("this message should have some buttons").Build()).
			SetComponents(api.NewRow(
				api.NewBlurpleButton("test", "test", api.NewEmoji("❌"), false),
				api.NewLinkButton("KittyBot", "https://kittybot.de", api.NewCustomEmoji("837665167780216852"), false),
				api.NewSelect("select", "placeholder", 1, 2,
					&api.SelectOption{
						Label:       "test1",
						Value:       "value1",
						Default:     false,
						Description: "value1",
					},
					&api.SelectOption{
						Label:       "test2",
						Value:       "value2",
						Default:     false,
						Description: "value2",
					},
				),
			)).
			Build(),
		); err != nil {
			logger.Errorf("error sending interaction response: %s", err)
		}

	case "addrole":
		user := event.Option("member").User()
		role := event.Option("role").Role()
		err := event.Disgo().RestClient().AddMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
		if err == nil {
			_ = event.Reply(api.NewCommandResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(green).SetDescriptionf("Added %s to %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewCommandResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(red).SetDescriptionf("Failed to add %s to %s", role, user).Build(),
			).Build())
		}

	case "removerole":
		user := event.Option("member").User()
		role := event.Option("role").Role()
		err := event.Disgo().RestClient().RemoveMemberRole(*event.Interaction.GuildID, user.ID, role.ID)
		if err == nil {
			_ = event.Reply(api.NewCommandResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(65280).SetDescriptionf("Removed %s from %s", role, user).Build(),
			).Build())
		} else {
			_ = event.Reply(api.NewCommandResponseBuilder().AddEmbeds(
				api.NewEmbedBuilder().SetColor(16711680).SetDescriptionf("Failed to remove %s from %s", role, user).Build(),
			).Build())
		}
	}
}

func messageListener(event *events.GuildMessageCreateEvent) {
	if event.Message.Author.IsBot {
		return
	}
	if event.Message.Content == nil {
		return
	}

	switch *event.Message.Content {
	case "ping":
		_, _ = event.Message.Reply(api.NewMessageBuilder().SetContent("pong").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "pong":
		_, _ = event.Message.Reply(api.NewMessageBuilder().SetContent("ping").SetAllowedMentions(&api.AllowedMentions{RepliedUser: false}).Build())

	case "dm":
		go func() {
			channel, err := event.Message.Author.OpenDMChannel()
			if err != nil {
				_ = event.Message.AddReaction("❌")
				return
			}
			_, err = channel.SendMessage(api.NewMessageBuilder().SetContent("helo").Build())
			if err == nil {
				_ = event.Message.AddReaction("✅")
			} else {
				_ = event.Message.AddReaction("❌")
			}
		}()
	}
}

func ptrBool(bool bool) *bool {
	return &bool
}
