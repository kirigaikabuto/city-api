package main

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kirigaikabuto/city-api/applications"
	"github.com/kirigaikabuto/city-api/auth"
	"github.com/kirigaikabuto/city-api/comments"
	"github.com/kirigaikabuto/city-api/common"
	"github.com/kirigaikabuto/city-api/events"
	"github.com/kirigaikabuto/city-api/feedback"
	"github.com/kirigaikabuto/city-api/mdw"
	"github.com/kirigaikabuto/city-api/news"
	"github.com/kirigaikabuto/city-api/user_events"
	"github.com/kirigaikabuto/city-api/users"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	configName              = "main"
	configPath              = "/config/"
	version                 = "0.0.1"
	redisHost               = ""
	redisPort               = ""
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
	port                    = ""
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
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	postgresDatabaseName = os.Getenv("POSTGRES_DB")
	postgresParams = os.Getenv("POSTGRES_PARAM")
	postgresPortStr := os.Getenv("POSTGRES_PORT")
	postgresPort, _ = strconv.Atoi(postgresPortStr)
	postgresHost = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("PORT")
	s3endpoint = viper.GetString("s3.primary.s3endpoint")
	s3bucket = viper.GetString("s3.primary.s3bucket")
	s3accessKey = viper.GetString("s3.primary.s3accessKey")
	s3secretKey = viper.GetString("s3.primary.s3secretKey")
	s3uploadedFilesBasePath = viper.GetString("s3.primary.s3uploadedFilesBasePath")
	s3region = viper.GetString("s3.primary.s3region")
	redisHost = os.Getenv("REDIS_HOST")
	redisPort = os.Getenv("REDIS_PORT")
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
	//tkn store
	tknStore, err := mdw.NewTokenStore(mdw.RedisConfig{
		Host: redisHost,
		Port: redisPort,
	})
	if err != nil {
		return err
	}
	mdw := mdw.NewMiddleware(tknStore)
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
	usersPostgreStore, err := users.NewPostgresUsersStore(cfg)
	if err != nil {
		return err
	}
	usersPostgreStore.Create(&users.User{
		FirstName:   "yerassyl",
		LastName:    "tleugazy",
		Username:    "admin",
		Password:    "admin",
		Email:       "tleugazy98@gmail.com",
		PhoneNumber: "12323",
		Gender:      "male",
		AccessType:  "admin",
	})
	applicationPostgreStore, err := applications.NewPostgresApplicationStore(cfg)
	if err != nil {
		return err
	}
	applicationService := applications.NewApplicationService(applicationPostgreStore, s3Uploader, usersPostgreStore)
	applicationHttpEndpoints := applications.NewHttpEndpoints(setdata_common.NewCommandHandler(applicationService))
	//events
	eventsPostgreStore, err := events.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	eventService := events.NewService(eventsPostgreStore, s3Uploader)
	eventsHttpEndpoints := events.NewHttpEndpoints(setdata_common.NewCommandHandler(eventService))

	r := gin.Default()
	//r.Use(apiKewMdw.MakeCorsMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://94.247.128.130", "http://chistiy-gorod.kz"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://94.247.128.130" || origin == "http://chistiy-gorod.kz"
		},
		MaxAge: 12 * time.Hour,
	}))
	authService := auth.NewService(usersPostgreStore, tknStore, s3Uploader)
	authHttpEndpoints := auth.NewHttpEndpoints(setdata_common.NewCommandHandler(authService))

	//feedback
	feedbackPostgreStore, err := feedback.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	feedbackService := feedback.NewService(feedbackPostgreStore)
	feedbackHttpEndpoints := feedback.NewHttpEndpoints(setdata_common.NewCommandHandler(feedbackService))

	//comments
	commentsPostgreStore, err := comments.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	commentsService := comments.NewService(commentsPostgreStore)
	commentsHttpEnpoints := comments.NewHttpEndpoints(setdata_common.NewCommandHandler(commentsService))

	//news
	newsPostgreStore, err := news.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	newsService := news.NewService(s3Uploader, newsPostgreStore)
	newsHttpEndpoints := news.NewHttpEndpoints(setdata_common.NewCommandHandler(newsService))

	//user events
	userEventsPostgreStore, err := user_events.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	userEventsService := user_events.NewService(userEventsPostgreStore, usersPostgreStore, eventsPostgreStore)
	userEventsHttpEndpoints := user_events.NewHttpEndpoints(setdata_common.NewCommandHandler(userEventsService))

	appGroup := r.Group("/application")
	{
		appGroup.POST("/", applicationHttpEndpoints.MakeCreateApplication())
		appGroup.POST("/create", mdw.MakeMiddleware(), applicationHttpEndpoints.MakeCreateApplicationWithAuth())
		appGroup.PUT("/file", mdw.MakeMiddleware(), applicationHttpEndpoints.MakeUploadApplicationFile())
		appGroup.PUT("/status", mdw.MakeMiddleware(), applicationHttpEndpoints.MakeUpdateStatus())
		appGroup.GET("/type", applicationHttpEndpoints.MakeListApplicationByType())
		appGroup.GET("/id", applicationHttpEndpoints.MakeGetApplicationById())
		appGroup.GET("/list", applicationHttpEndpoints.MakeListApplication())
		appGroup.GET("/my", mdw.MakeMiddleware(), applicationHttpEndpoints.MakeAuthorizedUserListApplications())
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/", mdw.MakeMiddleware(), eventsHttpEndpoints.MakeCreateEvent())
		eventGroup.GET("/", mdw.MakeMiddleware(), eventsHttpEndpoints.MakeListEvent())
		eventGroup.GET("/my", mdw.MakeMiddleware(), eventsHttpEndpoints.MakeListEventByUserId())
		eventGroup.PUT("/document", mdw.MakeMiddleware(), eventsHttpEndpoints.MakeUploadDocument())
	}
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHttpEndpoints.MakeLoginEndpoint())
		authGroup.POST("/register", authHttpEndpoints.MakeRegisterEndpoint())
		authGroup.GET("/profile", mdw.MakeMiddleware(), authHttpEndpoints.MakeGetProfileEndpoint())
		authGroup.PUT("/profile", mdw.MakeMiddleware(), authHttpEndpoints.MakeUpdateProfileEndpoint())
		authGroup.PUT("/avatar", mdw.MakeMiddleware(), authHttpEndpoints.MakeUploadAvatarEndpoint())
	}
	feedbackGroup := r.Group("/feedback")
	{
		feedbackGroup.POST("/", feedbackHttpEndpoints.MakeCreateFeedback())
		feedbackGroup.GET("/", feedbackHttpEndpoints.MakeListFeedback())
	}
	commentsGroup := r.Group("/comment")
	{
		commentsGroup.POST("/", mdw.MakeMiddleware(), commentsHttpEnpoints.MakeCreateEndpoint())
		commentsGroup.GET("/", mdw.MakeMiddleware(), commentsHttpEnpoints.MakeListEndpoint())
		commentsGroup.GET("/obj", commentsHttpEnpoints.MakeListByObjTypeEndpoint())
	}
	newsGroup := r.Group("/news")
	{
		newsGroup.POST("/", mdw.MakeMiddleware(), newsHttpEndpoints.MakeCreateNews())
		newsGroup.PUT("/", mdw.MakeMiddleware(), newsHttpEndpoints.MakeUpdateNews())
		newsGroup.PUT("/photo", mdw.MakeMiddleware(), newsHttpEndpoints.MakeUploadPhoto())
		newsGroup.GET("/", newsHttpEndpoints.MakeListNews())
		newsGroup.GET("/id", newsHttpEndpoints.MakeGetNewsById())
		newsGroup.GET("/my", mdw.MakeMiddleware(), newsHttpEndpoints.MakeGetNewsByAuthorId())
	}
	userEventsGroup := r.Group("/user-events")
	{
		userEventsGroup.POST("/", mdw.MakeMiddleware(), userEventsHttpEndpoints.MakeCreateUserEvent())
		userEventsGroup.GET("/userId", mdw.MakeMiddleware(), userEventsHttpEndpoints.MakeListByUserId())
		userEventsGroup.GET("/eventId", mdw.MakeMiddleware(), userEventsHttpEndpoints.MakeListByEventId())
		userEventsGroup.GET("/", mdw.MakeMiddleware(), userEventsHttpEndpoints.MakeListUserEvents())
		userEventsGroup.GET("/id", mdw.MakeMiddleware(), userEventsHttpEndpoints.MakeGetUserEventById())
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
