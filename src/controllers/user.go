package controllers

import (
	"jwt/src/helpers"
	"jwt/src/models"
	"jwt/src/services"
	"jwt/src/utils"
	"jwt/src/validations"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Paginatex struct {
	Meta map[string]interface{} `json:"meta"`
	Data interface{}            `json:"data"`
}

type UserController interface {
	All(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UploadAvatar(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) All(ctx *gin.Context) {
	var users []models.User = c.userService.All()
	res := utils.Response(true, "success", users, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Find(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)

	if err != nil {
		res := utils.ResponseError("No param id was found", utils.EmptyObj{}, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	var user models.User = c.userService.Find(id)
	if (user == models.User{}) {
		res := utils.ResponseError("Data not found", utils.EmptyObj{}, nil)
		ctx.JSON(http.StatusNotFound, res)
	} else {
		res := utils.Response(true, "Detail user", user, nil)
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *userController) Create(ctx *gin.Context) {
	var userValidation validations.UserCreateValidation
	errReq := ctx.ShouldBind(&userValidation)
	if errReq != nil {
		res := utils.ResponseError("Invalid request", utils.EmptyObj{}, errReq)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	user := c.userService.Create(userValidation)
	res := utils.Response(true, "success", user, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (c *userController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := utils.ResponseError("No param id was found", utils.EmptyObj{}, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var user models.User = c.userService.Find(id)
	if (user == models.User{}) {
		res := utils.ResponseError("Data not found", utils.EmptyObj{}, nil)
		ctx.JSON(http.StatusNotFound, res)
	}

	var userValidation validations.UserUpdateValidation
	userValidation.ID = id
	errReq := ctx.ShouldBind(&userValidation)

	if errReq != nil {
		res := utils.ResponseError("Invalid request", utils.EmptyObj{}, errReq)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result := c.userService.Update(userValidation)
	res := utils.Response(true, "success", result, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (c *userController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := utils.ResponseError("No param id was found", utils.EmptyObj{}, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var user models.User = c.userService.Find(id)
	if (user == models.User{}) {
		res := utils.ResponseError("Data not found", utils.EmptyObj{}, nil)
		ctx.JSON(http.StatusNotFound, res)
	}

	result := c.userService.Delete(user)
	res := utils.Response(true, "success", result, nil)
	ctx.JSON(http.StatusCreated, res)
}

func (c *userController) UploadAvatar(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := utils.ResponseError("No param id was found", utils.EmptyObj{}, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var user models.User = c.userService.Find(id)
	if (user == models.User{}) {
		res := utils.ResponseError("Data not found", utils.EmptyObj{}, nil)
		ctx.JSON(http.StatusNotFound, res)
	}

	message, err := helpers.UploadFile(ctx, "avatar")

	if err != nil {
		res := utils.ResponseError(message, utils.EmptyObj{}, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var userValidation validations.UserAvatarValidation
	userValidation.ID = id
	userValidation.Avatar = message

	errReq := ctx.ShouldBind(&userValidation)
	if errReq != nil {
		res := utils.ResponseError("Invalid request", utils.EmptyObj{}, errReq)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result := c.userService.UploadAvatar(userValidation)
	res := utils.Response(true, "success", result, nil)
	ctx.JSON(http.StatusCreated, res)
}
