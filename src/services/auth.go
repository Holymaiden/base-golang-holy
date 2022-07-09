package services

import (
	"jwt/src/helpers"
	"jwt/src/models"
	"jwt/src/repositorys"
	"jwt/src/validations"
	"log"

	"github.com/mashingan/smapping"
)

type AuthService interface {
	FindByEmail(email string) models.User
	CreateUser(user validations.UserCreateValidation) models.User
	VerifyCredentials(email string, password string) interface{}
}

type authService struct {
	userRepository repositorys.UserRepository
}

func NewAuthService(userRepository repositorys.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (services *authService) FindByEmail(email string) models.User {
	return services.userRepository.FindByEmail(email)
}

func (services *authService) CreateUser(user validations.UserCreateValidation) models.User {
	userToCreate := models.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := services.userRepository.Create(userToCreate)
	return res
}

func (service *authService) VerifyCredentials(email string, password string) interface{} {
	res := service.userRepository.VerifyCredentials(email)
	if v, ok := res.(models.User); ok {
		comparePassword := helpers.ComparePassword(v.Password, []byte(password))
		if email == v.Email && comparePassword {
			return res
		}
		return false
	}
	return false

}
