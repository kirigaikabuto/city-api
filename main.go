package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kirigaikabuto/city-api/application"
	"github.com/kirigaikabuto/city-api/common"
	"github.com/kirigaikabuto/city-api/events"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"log"
	"os"
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
	port                    = "8080"
	flags                   = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Usage:       "path to .env config file",
			Value:       "./config/",
			Destination: &configPath,
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
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	gin.SetMode(gin.ReleaseMode)
	cfg := common.PostgresConfig{
		Host:     postgresHost,
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabaseName,
		Params:   postgresParams,
	}
	//application
	applicationPostgreStore, err := application.NewPostgresApplicationStore(cfg)
	if err != nil {
		return err
	}
	applicationService := application.NewApplicationService(applicationPostgreStore)
	applicationHttpEndpoints := application.NewHttpEndpoints(setdata_common.NewCommandHandler(applicationService))
	//events
	eventsPostgreStore, err := events.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	eventService := events.NewService(eventsPostgreStore)
	eventsHttpEnpoints := events.NewHttpEndpoints(setdata_common.NewCommandHandler(eventService))
	r := gin.Default()

	authGroup := r.Group("/application")
	{
		authGroup.POST("/", applicationHttpEndpoints.MakeCreateApplication())
		authGroup.GET("/", applicationHttpEndpoints.MakeListApplication())
	}
	searchGroup := r.Group("/search")
	{
		searchGroup.GET("/street", applicationHttpEndpoints.MakeSearchPlace())
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/", eventsHttpEnpoints.MakeCreateEvent())
		eventGroup.GET("/", eventsHttpEnpoints.MakeListEvent())
	}
	log.Println("app is running on port:" + port)
	return r.Run()
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
		log.Fatal(err)
	}
}
