package di

import (
	"context"
	"log"

	"gorm.io/gorm"

	"gofiber-template/application/serviceimpl"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/external/google"
	"gofiber-template/infrastructure/external/openai"
	"gofiber-template/infrastructure/postgres"
	"gofiber-template/infrastructure/redis"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/oauth"
	"gofiber-template/pkg/scheduler"
)

type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB             *gorm.DB
	RedisClient    *redis.RedisClient
	R2Storage      storage.R2Storage
	EventScheduler scheduler.EventScheduler

	// External API Clients
	GoogleSearchClient    *google.SearchClient
	GooglePlacesClient    *google.PlacesClient
	GoogleYouTubeClient   *google.YouTubeClient
	GoogleTranslateClient *google.TranslateClient
	OpenAIClient          *openai.AIClient

	// Repositories
	UserRepository          repositories.UserRepository
	TaskRepository          repositories.TaskRepository
	FileRepository          repositories.FileRepository
	JobRepository           repositories.JobRepository
	FolderRepository        repositories.FolderRepository
	FolderItemRepository    repositories.FolderItemRepository
	FavoriteRepository      repositories.FavoriteRepository
	SearchHistoryRepository repositories.SearchHistoryRepository
	AIChatSessionRepository    repositories.AIChatSessionRepository
	AIChatMessageRepository    repositories.AIChatMessageRepository
	PlaceAIContentRepository   repositories.PlaceAIContentRepository

	// Services
	UserService     services.UserService
	TaskService     services.TaskService
	FileService     services.FileService
	JobService      services.JobService
	SearchService   services.SearchService
	AIService       services.AIService
	FolderService   services.FolderService
	FavoriteService services.FavoriteService
	UtilityService  services.UtilityService
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initInfrastructure(); err != nil {
		return err
	}

	if err := c.initRepositories(); err != nil {
		return err
	}

	if err := c.initServices(); err != nil {
		return err
	}

	if err := c.initScheduler(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.Config = cfg
	log.Println("✓ Configuration loaded")
	return nil
}

func (c *Container) initInfrastructure() error {
	// Initialize Database
	dbConfig := postgres.DatabaseConfig{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db
	log.Println("✓ Database connected")

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		return err
	}
	log.Println("✓ Database migrated")

	// Initialize Redis
	redisConfig := redis.RedisConfig{
		Host:     c.Config.Redis.Host,
		Port:     c.Config.Redis.Port,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	}
	c.RedisClient = redis.NewRedisClient(redisConfig)

	// Test Redis connection
	if err := c.RedisClient.Ping(context.Background()); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("✓ Redis connected")
	}

	// Initialize Cloudflare R2 Storage
	r2Config := storage.R2Config{
		AccountID:       c.Config.R2.AccountID,
		AccessKeyID:     c.Config.R2.AccessKeyID,
		SecretAccessKey: c.Config.R2.SecretAccessKey,
		Bucket:          c.Config.R2.Bucket,
		PublicURL:       c.Config.R2.PublicURL,
	}
	r2Storage, err := storage.NewR2Storage(r2Config)
	if err != nil {
		log.Printf("Warning: R2 Storage initialization failed: %v", err)
	} else {
		c.R2Storage = r2Storage
		log.Println("✓ Cloudflare R2 Storage initialized")
	}

	// Initialize External API Clients
	if err := c.initExternalClients(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initExternalClients() error {
	// Initialize Google API Clients
	apiKey := c.Config.Google.APIKey
	mapsAPIKey := c.Config.Google.MapsAPIKey
	if mapsAPIKey == "" {
		mapsAPIKey = apiKey // fallback to main API key
	}

	c.GoogleSearchClient = google.NewSearchClient(apiKey, c.Config.Google.SearchEngineID)
	c.GooglePlacesClient = google.NewPlacesClient(mapsAPIKey)
	c.GoogleYouTubeClient = google.NewYouTubeClient(apiKey)
	c.GoogleTranslateClient = google.NewTranslateClient(apiKey)
	log.Println("✓ Google API clients initialized")

	// Initialize OpenAI Client
	c.OpenAIClient = openai.NewAIClient(c.Config.OpenAI.APIKey, c.Config.OpenAI.Model)
	log.Println("✓ OpenAI client initialized")

	// Initialize OAuth Clients
	oauth.InitGoogleOAuth()
	oauth.InitLineOAuth()
	log.Println("✓ OAuth clients initialized")

	return nil
}

func (c *Container) initRepositories() error {
	c.UserRepository = postgres.NewUserRepository(c.DB)
	c.TaskRepository = postgres.NewTaskRepository(c.DB)
	c.FileRepository = postgres.NewFileRepository(c.DB)
	c.JobRepository = postgres.NewJobRepository(c.DB)

	// STOU Smart Tour repositories
	c.FolderRepository = postgres.NewFolderRepository(c.DB)
	c.FolderItemRepository = postgres.NewFolderItemRepository(c.DB)
	c.FavoriteRepository = postgres.NewFavoriteRepository(c.DB)
	c.SearchHistoryRepository = postgres.NewSearchHistoryRepository(c.DB)
	c.AIChatSessionRepository = postgres.NewAIChatSessionRepository(c.DB)
	c.AIChatMessageRepository = postgres.NewAIChatMessageRepository(c.DB)
	c.PlaceAIContentRepository = postgres.NewPlaceAIContentRepository(c.DB)

	log.Println("✓ Repositories initialized")
	return nil
}

func (c *Container) initServices() error {
	c.UserService = serviceimpl.NewUserService(c.UserRepository, c.Config.JWT.Secret, c.R2Storage, c.Config.R2.PublicURL)
	c.TaskService = serviceimpl.NewTaskService(c.TaskRepository, c.UserRepository)
	c.FileService = serviceimpl.NewFileService(c.FileRepository, c.UserRepository, c.R2Storage)

	// STOU Smart Tour services
	c.SearchService = serviceimpl.NewSearchService(
		c.SearchHistoryRepository,
		c.PlaceAIContentRepository,
		c.GoogleSearchClient,
		c.GooglePlacesClient,
		c.GoogleYouTubeClient,
		c.OpenAIClient,
		c.RedisClient.GetClient(),
	)

	c.AIService = serviceimpl.NewAIService(
		c.AIChatSessionRepository,
		c.AIChatMessageRepository,
		c.SearchHistoryRepository,
		c.OpenAIClient,
		c.GoogleSearchClient,
		c.RedisClient.GetClient(),
	)

	c.FolderService = serviceimpl.NewFolderService(
		c.FolderRepository,
		c.FolderItemRepository,
		c.R2Storage,
	)

	c.FavoriteService = serviceimpl.NewFavoriteService(c.FavoriteRepository)

	c.UtilityService = serviceimpl.NewUtilityService(
		c.GoogleTranslateClient,
		c.RedisClient.GetClient(),
		c.Config,
	)

	log.Println("✓ Services initialized")
	return nil
}

func (c *Container) initScheduler() error {
	c.EventScheduler = scheduler.NewEventScheduler()
	c.JobService = serviceimpl.NewJobService(c.JobRepository, c.EventScheduler)

	// Start the scheduler
	c.EventScheduler.Start()
	log.Println("✓ Event scheduler started")

	// Load and schedule existing active jobs
	ctx := context.Background()
	jobs, _, err := c.JobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		log.Printf("Warning: Failed to load existing jobs: %v", err)
		return nil
	}

	activeJobCount := 0
	for _, job := range jobs {
		if job.IsActive {
			err := c.EventScheduler.AddJob(job.ID.String(), job.CronExpr, func() {
				c.JobService.ExecuteJob(ctx, job)
			})
			if err != nil {
				log.Printf("Warning: Failed to schedule job %s: %v", job.Name, err)
			} else {
				activeJobCount++
			}
		}
	}

	if activeJobCount > 0 {
		log.Printf("✓ Scheduled %d active jobs", activeJobCount)
	}

	return nil
}

func (c *Container) Cleanup() error {
	log.Println("Starting cleanup...")

	// Stop scheduler
	if c.EventScheduler != nil {
		if c.EventScheduler.IsRunning() {
			c.EventScheduler.Stop()
			log.Println("✓ Event scheduler stopped")
		} else {
			log.Println("✓ Event scheduler was already stopped")
		}
	}

	// Close Redis connection
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			log.Printf("Warning: Failed to close Redis connection: %v", err)
		} else {
			log.Println("✓ Redis connection closed")
		}
	}

	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	log.Println("✓ Cleanup completed")
	return nil
}

func (c *Container) GetServices() (services.UserService, services.TaskService, services.FileService, services.JobService) {
	return c.UserService, c.TaskService, c.FileService, c.JobService
}

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetHandlerServices() *handlers.Services {
	return &handlers.Services{
		UserService:     c.UserService,
		TaskService:     c.TaskService,
		FileService:     c.FileService,
		JobService:      c.JobService,
		SearchService:   c.SearchService,
		AIService:       c.AIService,
		FolderService:   c.FolderService,
		FavoriteService: c.FavoriteService,
		UtilityService:  c.UtilityService,
	}
}