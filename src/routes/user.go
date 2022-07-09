package routes

import (
	"jwt/src/controllers"
	"jwt/src/helpers"
	"jwt/src/middlewares"
	"jwt/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	userService    services.UserService       = services.NewUserService(userRepository)
	userController controllers.UserController = controllers.NewUserController(userService, jwtService)
)

func userRoute(v1 *gin.RouterGroup) *gin.RouterGroup {
	users := v1.Group("/user", middlewares.AuthorizeJWT(jwtService))
	{
		// users.GET("", userController.Index)
		users.GET("", func(context *gin.Context) {
			code := http.StatusOK

			pagination := helpers.Pagination(context)

			response := services.PaginationUser(userRepository, context, pagination)

			if !response.Code {
				code = http.StatusBadRequest
			}

			context.JSON(code, response)
		})
		users.POST("/", userController.Create)
		users.GET("/:id", userController.Find)
		users.PUT("/:id", userController.Update)
		users.DELETE("/:id", userController.Delete)
		users.POST("/:id/avatar", userController.UploadAvatar)
	}
	return users
}
