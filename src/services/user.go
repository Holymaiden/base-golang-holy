package services

import (
	"fmt"
	"jwt/src/models"
	"jwt/src/repositorys"
	"jwt/src/requests"
	"jwt/src/validations"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type UserService interface {
	All() []models.User
	Find(id uint64) models.User
	Create(model validations.UserCreateValidation) models.User
	Update(model validations.UserUpdateValidation) models.User
	Delete(model models.User) models.User
	UploadAvatar(model validations.UserAvatarValidation) models.User
}

type userService struct {
	userRepository repositorys.UserRepository
}

func NewUserService(userRepo repositorys.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) All() []models.User {
	return service.userRepository.All()
}

func (service *userService) Find(id uint64) models.User {
	return service.userRepository.Find(id)
}

func (service *userService) Create(model validations.UserCreateValidation) models.User {
	user := models.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&model))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	return service.userRepository.Create(user)
}

func (service *userService) Update(model validations.UserUpdateValidation) models.User {
	user := models.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&model))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	return service.userRepository.Update(user)
}

func (service *userService) Delete(model models.User) models.User {
	return service.userRepository.Delete(model)
}

func (service *userService) UploadAvatar(model validations.UserAvatarValidation) models.User {
	user := service.userRepository.Find(model.ID)
	user.Avatar = model.Avatar
	return service.userRepository.Update(user)
}

func PaginationUser(repository repositorys.UserRepository, context *gin.Context, pagination *requests.Pagination) requests.Response {

	operationResult, totalPages := repository.Pagination(pagination)

	if operationResult.Error != nil {
		return requests.Response{Code: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*requests.Pagination)

	// get current url path
	urlPath := context.Request.URL.Path

	// search query params
	searchQueryParams := ""

	for _, search := range pagination.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	// set first & last page pagination response
	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 1, pagination.Sort) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, totalPages, pagination.Sort) + searchQueryParams

	if data.Page > 1 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page, pagination.Sort) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page+1, pagination.Sort) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return requests.Response{Code: true, Data: data}
}
