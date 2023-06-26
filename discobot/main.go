package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/stablecog/sc-go/database"
	"github.com/stablecog/sc-go/database/repository"
	"github.com/stablecog/sc-go/discobot/commands"
	"github.com/stablecog/sc-go/log"
	"github.com/stablecog/sc-go/utils"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
)

func main() {
	// Load .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Warn("Error loading .env file (this is fine)", "err", err)
	}

	ctx := context.Background()

	// Setup redis
	redis, err := database.NewRedis(ctx)
	if err != nil {
		log.Fatal("Error connecting to redis", "err", err)
		os.Exit(1)
	}

	// Setup sql
	log.Info("🏡 Connecting to database...")
	dbconn, err := database.GetSqlDbConn(false)
	if err != nil {
		log.Fatal("Failed to connect to database", "err", err)
		os.Exit(1)
	}
	entClient, err := database.NewEntClient(dbconn)
	if err != nil {
		log.Fatal("Failed to create ent client", "err", err)
		os.Exit(1)
	}
	defer entClient.Close()
	// Run migrations
	// We can't run on supabase, :(
	if utils.GetEnv("RUN_MIGRATIONS", "") == "true" {
		log.Info("🦋 Running migrations...")
		if err := entClient.Schema.Create(ctx); err != nil {
			log.Fatal("Failed to run migrations", "err", err)
			os.Exit(1)
		}
	}

	// Create repository (database access)
	repo := &repository.Repository{
		DB:       entClient,
		ConnInfo: dbconn,
		Redis:    redis,
		Ctx:      ctx,
	}

	// Setup commands
	cmdWrapper := commands.DiscordCommands{
		Repo: repo,
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Info("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, 1)
	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, cmdWrapper.AuthenticateCommand())
	if err != nil {
		log.Fatalf("Cannot create '%v' command: %v", cmdWrapper.AuthenticateCommand().Name, err)
	}
	registeredCommands[0] = cmd

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmdWrapper.AuthenticateHandler(s, i)
	})

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Info("Removing commands...")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Fatalf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Info("Gracefully shutting down.")
}
