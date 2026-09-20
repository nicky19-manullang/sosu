package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sosu-backend/internal/fetcher"
	"sosu-backend/internal/handler"
	"sosu-backend/internal/middleware"
	"sosu-backend/internal/model"
	"sosu-backend/internal/repository"
	"sosu-backend/internal/scheduler"
	"sosu-backend/internal/service"
)

func runMigrations(db *gorm.DB) {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal("failed to read migrations folder: ", err)
	}
	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatal("failed to read migration file: ", f, err)
		}
		if err := db.Exec(string(content)).Error; err != nil {
			log.Fatal("failed to run migration: ", f, err)
		}
		log.Println("migration applied:", f)
	}
}

// allowedOrigins membaca ALLOWED_ORIGINS (dipisah koma) dari environment.
// Kalau kosong, default ke origin dev lokal.
func allowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		return []string{"http://localhost:5173"}
	}
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

// adminOnly melindungi endpoint operasional (fetch manual, tes digest email).
// Endpoint hanya aktif kalau ADMIN_TOKEN diisi di environment; klien harus
// mengirim header X-Admin-Token. Kalau ADMIN_TOKEN kosong, endpoint dianggap tidak ada (404).
func adminOnly() gin.HandlerFunc {
	token := os.Getenv("ADMIN_TOKEN")
	return func(c *gin.Context) {
		if token == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		got := c.GetHeader("X-Admin-Token")
		if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from system env")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")

	if err := db.AutoMigrate(&model.Hotspot{}, &model.AQIReading{}, &model.User{}, &model.Subscription{}); err != nil {
		log.Fatal("failed to auto-migrate: ", err)
	}

	runMigrations(db)

	hotspotRepo := repository.NewHotspotRepository(db)
	aqiRepo := repository.NewAQIRepository(db)

	firmsFetcher := fetcher.NewFIRMSFetcher(os.Getenv("FIRMS_MAP_KEY"), hotspotRepo)
	openaqFetcher := fetcher.NewOpenAQFetcher(os.Getenv("OPENAQ_API_KEY"), aqiRepo)

	statusService := service.NewStatusService(hotspotRepo, aqiRepo)
	trendService := service.NewTrendService(hotspotRepo, aqiRepo)
	trendHandler := handler.NewTrendHandler(trendService)
	statusHandler := handler.NewStatusHandler(statusService)
	hotspotHandler := handler.NewHotspotHandler(hotspotRepo)
	aqiHandler := handler.NewAQIHandler(aqiRepo)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(userRepo, authService)

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	emailService := service.NewEmailService()
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionRepo, userRepo, emailService, statusService)

	sched := scheduler.New(firmsFetcher, openaqFetcher, statusService, subscriptionRepo, userRepo, emailService, hotspotRepo)
	go sched.RunInitialFetch()
	sched.Start()
	defer sched.Stop()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(500, gin.H{"status": "db not connected"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Endpoint operasional: hanya bisa diakses dengan header X-Admin-Token.
	admin := router.Group("/api/v1", adminOnly())
	{
		admin.POST("/fetch/hotspots", func(c *gin.Context) {
			if err := firmsFetcher.Run(); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "fetched"})
		})

		admin.POST("/fetch/aqi", func(c *gin.Context) {
			if err := openaqFetcher.Run(); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "fetched"})
		})

		admin.POST("/fetch/hotspots/global", func(c *gin.Context) {
			if err := firmsFetcher.RunGlobal(); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "fetched"})
		})

		admin.POST("/test/weekly-digest", func(c *gin.Context) {
			sched.TriggerDigestNow()
			c.JSON(200, gin.H{"status": "digest triggered"})
		})
	}

	api := router.Group("/api/v1")
	{
		api.GET("/status", statusHandler.GetStatus)
		api.GET("/hotspots", hotspotHandler.GetHotspots)
		api.GET("/hotspots/bbox", hotspotHandler.GetHotspotsInBBox)
		api.GET("/trend", trendHandler.GetTrend)

		api.GET("/aqi/nearest", aqiHandler.GetNearest)

		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		authorized := api.Group("/")
		authorized.Use(middleware.AuthRequired(authService))
		{
			authorized.POST("/subscribe", subscriptionHandler.Subscribe)
			authorized.POST("/unsubscribe", subscriptionHandler.Unsubscribe)
			authorized.GET("/my-subscriptions", subscriptionHandler.MySubscriptions)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}