package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kirigaikabuto/city-api/api_keys"
	"github.com/kirigaikabuto/city-api/applications"
	"github.com/kirigaikabuto/city-api/auth"
	"github.com/kirigaikabuto/city-api/common"
	"github.com/kirigaikabuto/city-api/events"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configName              = "main"
	configPath              = "/config/"
	version                 = "0.0.1"
	s3endpoint              = ""
	s3bucket                = ""
	s3accessKey             = ""
	s3secretKey             = ""
	s3uploadedFilesBasePath = ""
	s3region                = ""
	postgresUser            = ""
	postgresPassword        = ""
	postgresDatabaseName    = ""
	postgresHost            = ""
	postgresPort            = 5432
	postgresParams          = ""
	port                    = "5000"
	flags                   = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Usage:       "path to .env config file",
			Value:       "main",
			Destination: &configName,
		},
	}
)

func parseEnvFile() {
	filepath, err := os.Getwd()
	if err != nil {
		panic("main, get rootDir error" + err.Error())
		return
	}
	viper.AddConfigPath(filepath + configPath)
	viper.SetConfigName(configName)
	err = viper.ReadInConfig()
	if err != nil {
		panic("main, fatal error while reading config file: " + err.Error())
		return
	}
	postgresUser = viper.GetString("db.primary.user")
	postgresPassword = viper.GetString("db.primary.pass")
	postgresDatabaseName = viper.GetString("db.primary.name")
	postgresParams = viper.GetString("db.primary.param")
	postgresPort = viper.GetInt("db.primary.port")
	postgresHost = viper.GetString("db.primary.host")
	s3endpoint = viper.GetString("s3.primary.s3endpoint")
	s3bucket = viper.GetString("s3.primary.s3bucket")
	s3accessKey = viper.GetString("s3.primary.s3accessKey")
	s3secretKey = viper.GetString("s3.primary.s3secretKey")
	s3uploadedFilesBasePath = viper.GetString("s3.primary.s3uploadedFilesBasePath")
	s3region = viper.GetString("s3.primary.s3region")
}

func run(c *cli.Context) error {
	parseEnvFile()
	gin.SetMode(gin.ReleaseMode)
	cfg := common.PostgresConfig{
		Host:     postgresHost,
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabaseName,
		Params:   postgresParams,
	}
	//applications
	s3Uploader, err := common.NewS3Uploader(
		s3endpoint,
		s3accessKey,
		s3secretKey,
		s3bucket,
		s3uploadedFilesBasePath,
		s3region)
	if err != nil {
		return err
	}
	applicationPostgreStore, err := applications.NewPostgresApplicationStore(cfg)
	if err != nil {
		return err
	}
	applicationService := applications.NewApplicationService(applicationPostgreStore, s3Uploader)
	applicationHttpEndpoints := applications.NewHttpEndpoints(setdata_common.NewCommandHandler(applicationService))
	//events
	eventsPostgreStore, err := events.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	eventService := events.NewService(eventsPostgreStore)
	eventsHttpEndpoints := events.NewHttpEndpoints(setdata_common.NewCommandHandler(eventService))

	apiKeyStore, err := api_keys.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	apiKeyHttpEndpoints := api_keys.NewHttpEndpoints(setdata_common.NewCommandHandler(apiKeyStore))
	apiKewMdw := auth.NewApiKeyMdw(apiKeyStore)
	r := gin.Default()
	//r.Use(apiKewMdw.MakeCorsMiddleware())
	appGroup := r.Group("/application")
	{
		appGroup.POST("/", applicationHttpEndpoints.MakeCreateApplication())
		appGroup.PUT("/file", applicationHttpEndpoints.MakeUploadApplicationFile())
		appGroup.PUT("/status", applicationHttpEndpoints.MakeUpdateStatus())
		appGroup.GET("/type", applicationHttpEndpoints.MakeListApplicationByType())
		appGroup.GET("/id", applicationHttpEndpoints.MakeGetApplicationById())
		appGroup.GET("/list", applicationHttpEndpoints.MakeListApplication())
	}
	searchGroup := r.Group("/search")
	{
		searchGroup.GET("/street", apiKewMdw.MakeApiKeyMiddleware(), applicationHttpEndpoints.MakeSearchPlace())
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/", apiKewMdw.MakeApiKeyMiddleware(), eventsHttpEndpoints.MakeCreateEvent())
		eventGroup.GET("/", apiKewMdw.MakeApiKeyMiddleware(), eventsHttpEndpoints.MakeListEvent())
	}
	apiKeyGroup := r.Group("/api-key")
	{
		apiKeyGroup.POST("/", apiKeyHttpEndpoints.MakeCreateApiKey())
		apiKeyGroup.GET("/", apiKeyHttpEndpoints.MakeListApiKey())
	}
	log.Info().Msg("app is running on port:" + port)
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Server ListenAndServe error")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting.")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "city api"
	app.Description = ""
	app.Usage = "city api"
	app.UsageText = "city api"
	app.Version = version
	app.Flags = flags
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		log.Err(err)
	}
}
