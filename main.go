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
	file_storage "github.com/kirigaikabuto/city-api/file-storage"
	"github.com/kirigaikabuto/city-api/mdw"
	"github.com/kirigaikabuto/city-api/news"
	sms_store "github.com/kirigaikabuto/city-api/sms-store"
	"github.com/kirigaikabuto/city-api/user_events"
	"github.com/kirigaikabuto/city-api/users"
	setdata_common "github.com/kirigaikabuto/setdata-common"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
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
	pulseEmailFrom          = ""
	pulseClientId           = ""
	pulseClientSecret       = ""
	pulseBasicUrl           = ""
)

func parseEnvFile() {
	postgresUser = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	postgresDatabaseName = os.Getenv("POSTGRES_DB")
	postgresParams = os.Getenv("POSTGRES_PARAM")
	postgresPortStr := os.Getenv("POSTGRES_PORT")
	postgresPort, _ = strconv.Atoi(postgresPortStr)
	postgresHost = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("PORT")
	s3endpoint = os.Getenv("S3_ENDPOINT")
	s3bucket = os.Getenv("S3_BUCKET")
	s3accessKey = os.Getenv("S3_ACCESS_KEY")
	s3secretKey = os.Getenv("S3_SECRET_KEY")
	s3uploadedFilesBasePath = os.Getenv("S3_FILE_UPLOAD_PATH")
	s3region = os.Getenv("S3_REGION")
	redisHost = os.Getenv("REDIS_HOST")
	redisPort = os.Getenv("REDIS_PORT")
	pulseBasicUrl = os.Getenv("PULSE_BASIC_URL")
	pulseClientId = os.Getenv("PULSE_CLIENT_ID")
	pulseClientSecret = os.Getenv("PULSE_CLIENT_SECRET")
	pulseEmailFrom = os.Getenv("EMAIL_FROM")
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
	mdwEndpoint := mdw.NewMiddleware(tknStore)
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
		IsVerified:  true,
	})
	fileStoragePostgresStore, err := file_storage.NewPostgresStore(cfg)
	if err != nil {
		return err
	}
	applicationPostgreStore, err := applications.NewPostgresApplicationStore(cfg, fileStoragePostgresStore)
	if err != nil {
		return err
	}
	applicationService := applications.NewApplicationService(applicationPostgreStore, s3Uploader, usersPostgreStore, fileStoragePostgresStore)
	applicationHttpEndpoints := applications.NewHttpEndpoints(setdata_common.NewCommandHandler(applicationService))
	//events
	eventsPostgreStore, err := events.NewPostgresStore(cfg, fileStoragePostgresStore)
	if err != nil {
		return err
	}
	eventService := events.NewService(eventsPostgreStore, s3Uploader, fileStoragePostgresStore, usersPostgreStore)
	eventsHttpEndpoints := events.NewHttpEndpoints(setdata_common.NewCommandHandler(eventService))
	//email
	emailStore := sms_store.NewPulseEmailStore(common.PulseEmailConfig{
		EmailFrom:    pulseEmailFrom,
		ClientId:     pulseClientId,
		ClientSecret: pulseClientSecret,
		BasicUrl:     pulseBasicUrl,
	})

	r := gin.Default()
	//r.Use(apiKewMdw.MakeCorsMiddleware())
	allowOrigins := []string{"http://94.247.128.130",
		"http://chistyi-gorod.kz",
		"http://37.99.44.126",
		"http://172.19.0.1",
		"http://localhost:63342",
		"*"}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			for _, v := range allowOrigins {
				if v == origin {
					return true
				}
			}
			return true
		},
		MaxAge: 72 * time.Hour,
	}))
	authService := auth.NewService(usersPostgreStore, tknStore, s3Uploader, emailStore)
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
	commentsService := comments.NewService(commentsPostgreStore, eventsPostgreStore, applicationPostgreStore, usersPostgreStore)
	commentsHttpEndpoints := comments.NewHttpEndpoints(setdata_common.NewCommandHandler(commentsService))

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
		appGroup.POST("/create", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeCreateApplicationWithAuth())
		appGroup.PUT("/file", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeUploadApplicationFile())
		appGroup.PUT("/status", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeUpdateStatus())
		appGroup.GET("/type", applicationHttpEndpoints.MakeListApplicationByType())
		appGroup.GET("/id", applicationHttpEndpoints.MakeGetApplicationById())
		appGroup.GET("/list", applicationHttpEndpoints.MakeListApplication())
		appGroup.GET("/my", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeAuthorizedUserListApplications())
		appGroup.PUT("/", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeUpdateApplication())
		appGroup.DELETE("/", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeRemoveApplication())
		appGroup.GET("/list/address-auth", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeListByAddressWithAuth())
		appGroup.GET("/list/address", applicationHttpEndpoints.MakeListByAddress())
		appGroup.PUT("/multiple/file", mdwEndpoint.MakeMiddleware(), applicationHttpEndpoints.MakeUploadMultipleFiles())
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/", mdwEndpoint.MakeMiddleware(), eventsHttpEndpoints.MakeCreateEvent())
		eventGroup.GET("/", eventsHttpEndpoints.MakeListEvent())
		eventGroup.GET("/my", mdwEndpoint.MakeMiddleware(), eventsHttpEndpoints.MakeListEventByUserId())
		eventGroup.PUT("/document", mdwEndpoint.MakeMiddleware(), eventsHttpEndpoints.MakeUploadDocument())
		eventGroup.GET("/id", eventsHttpEndpoints.MakeGetEventById())
		eventGroup.PUT("/multiple/file", mdwEndpoint.MakeMiddleware(), eventsHttpEndpoints.MakeUploadMultipleFiles())
	}
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHttpEndpoints.MakeLoginEndpoint())
		authGroup.POST("/register", authHttpEndpoints.MakeRegisterEndpoint())
		authGroup.GET("/profile", mdwEndpoint.MakeMiddleware(), authHttpEndpoints.MakeGetProfileEndpoint())
		authGroup.PUT("/profile", mdwEndpoint.MakeMiddleware(), authHttpEndpoints.MakeUpdateProfileEndpoint())
		authGroup.PUT("/avatar", mdwEndpoint.MakeMiddleware(), authHttpEndpoints.MakeUploadAvatarEndpoint())
		authGroup.GET("/verify", authHttpEndpoints.MakeVerifyCodeEndpoint())
		authGroup.PUT("/reset-password-request", authHttpEndpoints.MakeResetPasswordRequestEndpoint())
		authGroup.PUT("/reset-password", authHttpEndpoints.MakeResetPasswordEndpoint())
		authGroup.DELETE("/remove", mdwEndpoint.MakeMiddleware(), authHttpEndpoints.MakeRemoveAccount())
	}
	feedbackGroup := r.Group("/feedback")
	{
		feedbackGroup.POST("/", feedbackHttpEndpoints.MakeCreateFeedback())
		feedbackGroup.GET("/", feedbackHttpEndpoints.MakeListFeedback())
	}
	commentsGroup := r.Group("/comment")
	{
		commentsGroup.POST("/", mdwEndpoint.MakeMiddleware(), commentsHttpEndpoints.MakeCreateEndpoint())
		commentsGroup.GET("/", mdwEndpoint.MakeMiddleware(), commentsHttpEndpoints.MakeListEndpoint())
		commentsGroup.GET("/obj", commentsHttpEndpoints.MakeListByObjTypeEndpoint())
		commentsGroup.GET("/objId", commentsHttpEndpoints.MakeListByObjectId())
	}
	newsGroup := r.Group("/news")
	{
		newsGroup.POST("/", mdwEndpoint.MakeMiddleware(), newsHttpEndpoints.MakeCreateNews())
		newsGroup.PUT("/", mdwEndpoint.MakeMiddleware(), newsHttpEndpoints.MakeUpdateNews())
		newsGroup.PUT("/photo", mdwEndpoint.MakeMiddleware(), newsHttpEndpoints.MakeUploadPhoto())
		newsGroup.GET("/", newsHttpEndpoints.MakeListNews())
		newsGroup.GET("/id", newsHttpEndpoints.MakeGetNewsById())
		newsGroup.GET("/my", mdwEndpoint.MakeMiddleware(), newsHttpEndpoints.MakeGetNewsByAuthorId())
	}
	userEventsGroup := r.Group("/user-events")
	{
		userEventsGroup.POST("/", mdwEndpoint.MakeMiddleware(), userEventsHttpEndpoints.MakeCreateUserEvent())
		userEventsGroup.GET("/userId", mdwEndpoint.MakeMiddleware(), userEventsHttpEndpoints.MakeListByUserId())
		userEventsGroup.GET("/eventId", mdwEndpoint.MakeMiddleware(), userEventsHttpEndpoints.MakeListByEventId())
		userEventsGroup.GET("/", mdwEndpoint.MakeMiddleware(), userEventsHttpEndpoints.MakeListUserEvents())
		userEventsGroup.GET("/id", mdwEndpoint.MakeMiddleware(), userEventsHttpEndpoints.MakeGetUserEventById())
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
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		log.Err(err)
	}
}
