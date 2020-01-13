package main

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/commands"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discordToken = "***REMOVED***"

const interval = 5 * time.Second
const maxColour = 16777216

func main() {
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		panic(fmt.Errorf("error creating Discord session: %w", err))
	}

	err = dg.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection: %w", err))
	}
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	guilds, err := dg.UserGuilds(0, "", "")
	if err != nil {
		panic(fmt.Errorf("error getting user guilds: %w", err))
	}

	guildRoles, err := guildroles.New(dg, guilds)
	if err != nil {
		panic(fmt.Errorf("error finding/creating rainbow roles: %w", err))
	}

	rand.Seed(time.Now().Unix())

	commands.Setup(dg, guildRoles)

	timer := time.NewTicker(interval)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-timer.C:
			err := changeRoleColours(dg, guildRoles)
			if err != nil {
				fmt.Println(err)
			}
		case <-sc:
			fmt.Println("Shutting down")
			return
		}
	}
}

func changeRoleColours(s *discordgo.Session, guildRoles guildroles.GuildRoles) error {
	for _, guildRole := range guildRoles {
		colour := rand.Intn(maxColour)

		_, err := s.GuildRoleEdit(guildRole.GuildId, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
		if err != nil {
			return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildId, err)
		}
	}

	return nil
}
