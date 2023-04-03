// Sets up a CLI to trigger the various cron jobs
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/stablecog/sc-go/cron/discord"
	"github.com/stablecog/sc-go/cron/jobs"
	"github.com/stablecog/sc-go/database"
	"github.com/stablecog/sc-go/database/repository"
	"github.com/stablecog/sc-go/log"
	"github.com/stablecog/sc-go/server/analytics"
	"github.com/stablecog/sc-go/utils"
	stripe "github.com/stripe/stripe-go/v74/client"
)

var Version = "dev"
var CommitMsg = "dev"

func usage() {
	fmt.Printf("Usage %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	return
}

func main() {
	log.Infof("SC Cron %s", Version)
	showHelp := flag.Bool("help", false, "Show help")
	healthCheck := flag.Bool("healthCheck", false, "Run the health check job")
	syncMeili := flag.Bool("syncMeili", false, "Sync the meili index")
	stats := flag.Bool("stats", false, "Run the stats job")
	allJobs := flag.Bool("all", false, "Run all jobs in a blocking process")
	flag.Parse()

	if *showHelp {
		usage()
		os.Exit(0)
	}

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

	// Create repostiory
	// Create repository (database access)
	repo := &repository.Repository{
		DB:       entClient,
		ConnInfo: dbconn,
		Redis:    redis,
		Ctx:      ctx,
	}

	// Create analytics service
	analyticsService := analytics.NewAnalyticsService()
	defer analyticsService.Close()

	// Create stripe client
	stripeClient := stripe.New(utils.GetEnv("STRIPE_SECRET_KEY", ""), nil)

	// Create a job runner
	jobRunner := jobs.JobRunner{
		Repo:    repo,
		Redis:   redis,
		Ctx:     ctx,
		Meili:   database.NewMeiliSearchClient(),
		Discord: discord.NewDiscordHealthTracker(ctx),
		Track:   analyticsService,
		Stripe:  stripeClient,
	}

	if *healthCheck {
		err := jobRunner.CheckHealth(jobs.NewJobLogger("HEALTH"))
		if err != nil {
			log.Fatal("Error running health check", "err", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *syncMeili {
		err := jobRunner.SyncMeili(jobs.NewJobLogger("MEILI_SYNC"))
		if err != nil {
			log.Fatal("Error syncing meili", "err", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *stats {
		err := jobRunner.GetAndSetStats(jobs.NewJobLogger("STATS"))
		if err != nil {
			log.Fatal("Error running stats job", "err", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *allJobs {
		// Get models, schedulers and put in cache
		log.Info("📦 Populating cache...")
		err = repo.UpdateCache()
		if err != nil {
			// ! Not getting these is fatal and will result in crash
			panic(err)
		}
		log.Info("🏡 Starting all jobs...")
		s := gocron.NewScheduler(time.UTC)
		s.Every(60).Seconds().Do(jobRunner.SyncMeili, jobs.NewJobLogger("MEILI_SYNC"))
		s.Every(60).Seconds().Do(jobRunner.GetAndSetStats, jobs.NewJobLogger("STATS"))
		if utils.GetEnv("DISCORD_WEBHOOK_URL", "") != "" {
			s.Every(60).Seconds().Do(jobRunner.CheckHealth, jobs.NewJobLogger("HEALTH"))
		}
		s.Every(60).Seconds().Do(jobRunner.AddFreeCreditsToEligibleUsers, jobs.NewJobLogger("FREE_CREDITS"))
		// Sync stripe
		s.Every(10).Minutes().Do(jobRunner.SyncStripe, jobs.NewJobLogger("STRIPE_SYNC"))
		// cache update
		s.Every(5).Minutes().StartAt(time.Now().Add(5 * time.Minute)).Do(func() {
			log.Info("📦 Updating cache...")
			err = repo.UpdateCache()
			if err != nil {
				log.Error("Error updating cache", "err", err)
			}
		})
		// Auto upscale
		// go jobRunner.StartAutoUpscaleJob(jobs.NewJobLogger("AUTO_UPSCALE"))
		s.StartBlocking()
		os.Exit(0)
	}

	// Generic path, they didn't say what they wanted.
	usage()
	os.Exit(1)

}
