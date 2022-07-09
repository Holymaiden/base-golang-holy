package controllers

import (
	"fmt"
	"jwt/src/models"
	"jwt/src/requests"
	"jwt/src/services"
	"jwt/src/utils"
	"jwt/src/validations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var validation validations.LoginValidation
	fmt.Print(validation)
	if errVal := ctx.ShouldBind(&validation); errVal != nil {
		res := utils.ResponseError("Invalid request", nil, errVal)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result := c.authService.VerifyCredentials(validation.Email, validation.Password)
	if user, ok := result.(models.User); ok {
		token := c.jwtService.GenerateToken(user.Email)
		tokenResponse := requests.LoginResponse{
			ID:          user.ID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			AccessToken: token,
		}
		res := utils.Response(true, "Login success", gin.H{"token": tokenResponse}, nil)
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := utils.Response(false, "Please check again your credential", "Invalid Credential", nil)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}

func (c *authController) Register(ctx *gin.Context) {

}

func (c *authController) RefreshToken(ctx *gin.Context) {

}
