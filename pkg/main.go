package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	pkg "git.solsynth.dev/matrix/nucleus/pkg/internal"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/gap"
	"github.com/fatih/color"

	"git.solsynth.dev/matrix/nucleus/pkg/internal/cache"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/database"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/grpc"

	"git.solsynth.dev/matrix/nucleus/pkg/internal/server"
	"git.solsynth.dev/matrix/nucleus/pkg/internal/services"
	"github.com/robfig/cron/v3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Booting screen
	fmt.Println(color.YellowString(" __  __       _        _\n|  \\/  | __ _| |_ _ __(_)_  __\n| |\\/| |/ _` | __| '__| \\ \\/ /\n| |  | | (_| | |_| |  | |>  <\n|_|  |_|\\__,_|\\__|_|  |_/_/\\_\\"))
	fmt.Printf("%s v%s\n", color.New(color.FgHiYellow).Add(color.Bold).Sprintf("Matrix.Nucleus"), pkg.AppVersion)
	fmt.Printf("The server side software of Matrix Software Marketplace\n")
	color.HiBlack("=====================================================\n")

	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Connect to nexus
	if err := gap.InitializeToNexus(); err != nil {
		log.Error().Err(err).Msg("An error occurred when registering service to nexus...")
	}

	// Load keypair
	if reader, err := sec.NewInternalTokenReader(viper.GetString("security.internal_public_key")); err != nil {
		log.Error().Err(err).Msg("An error occurred when reading internal public key for jwt. Authentication related features will be disabled.")
	} else {
		server.IReader = reader
		log.Info().Msg("Internal jwt public key loaded.")
	}

	// Connect to database
	if err := database.NewGorm(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}

	// Initialize cache
	if err := cache.NewStore(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when initializing cache.")
	}

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 60m", services.DoAutoDatabaseCleanup)
	quartz.Start()

	// Server
	go server.NewServer().Listen()

	// Grpc Server
	go grpc.NewGrpc().Listen()

	// Messages
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	quartz.Stop()
}
