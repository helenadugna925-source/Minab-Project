package main

import (
	"log"
	
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"local-event-backend/handlers"
	"local-event-backend/utils"
)

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, relying on system env variables")
	}

	// 2. Connect DB
	_, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// 3. Initialize Uploads Folder
	utils.InitUploadPath() 

	// 4. Create Gin
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 5. Improved CORS Middleware for Hasura & Nuxt
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow Hasura Docker
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-hasura-admin-secret")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 6. Serve Uploaded Images
	r.Static("/uploads", "./uploads") 

	// 7. Register Routes
	r.POST("/login", handlers.LoginHandler)
	r.POST("/signup", handlers.SignupHandler)
	r.POST("/upload", handlers.UploadHandler)
	r.POST("/create-event", handlers.CreateEventHandler)
	
	// CHAPA ROUTES
	r.POST("/initialize-payment", handlers.InitializePaymentHandler) 
	r.POST("/webhook/chapa", handlers.ChapaWebhookHandler)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 8. Start Server on a STABLE port for Hasura
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Keep it 8082 so Hasura always knows where to find it
	}
	addr := ":" + port

	log.Printf("🚀 Server running on http://localhost:%s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server exited: %v", err)
	}
}