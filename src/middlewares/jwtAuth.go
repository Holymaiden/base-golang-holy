package middlewares

import (
	"jwt/src/services"
	"jwt/src/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := utils.ResponseError("Failed to process request", utils.EmptyObj{}, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		extractedToken := strings.Split(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(extractedToken[1])
		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["id"])
			log.Println("Claim[issuer] :", claims["iss"])
		} else {
			log.Println(err)
			response := utils.ResponseError("Token is not valid", utils.EmptyObj{}, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
