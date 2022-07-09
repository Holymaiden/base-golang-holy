package routes

import (
	"jwt/config"
	"jwt/src/controllers"
	"jwt/src/repositorys"
	"jwt/src/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	db             *gorm.DB                   = config.SetupConnection()
	userRepository repositorys.UserRepository = repositorys.NewUserRepository(db)
	jwtService     services.JWTService        = services.NewJWTService()

	authServices   services.AuthService       = services.NewAuthService(userRepository)
	authController controllers.AuthController = controllers.NewAuthController(authServices, jwtService)
)

func SetupRouter() *gin.Engine {
	// defer config.CloseConnection(db)

	// Gin instance
	r := gin.Default()
	// if !envConfig.Debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// Routes
	v1 := r.Group("api/v1")
	{
		auth := v1.Group("auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}
		userRoute(v1)
	}
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			//c.Next()
			return
		}
		c.Next()
	}
}
