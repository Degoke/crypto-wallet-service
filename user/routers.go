package user

import (
	"errors"
	"net/http"

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/gin-gonic/gin"
)

func RegisterUnprotectedRoutes(router *gin.RouterGroup) {
	router.POST("/register", Register)
	router.POST("/login", Login)
}

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/me", GetUser)
}

func Register (c *gin.Context) {
	validator := NewUserValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError(err))
		return
	}

	_, err := FindOne(&User{Email: validator.User.Email})

	if err == nil {
		c.JSON(http.StatusConflict, common.NewError("registration", errors.New("email already exists")))
		return
	}

	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, common.NewError("registration", err))
		return
	}

	userModel := validator.userModel
	if err := Save(&userModel); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("registration", err))
		return
	}
	c.Set("user", userModel)

	serializer := UserSerializer{c}

	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}

func Login (c *gin.Context) {
	validator := NewLoginValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError(err))
		return
	}

	userModel, err := FindOne(&User{Email: validator.User.Email})

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("Unauthorized", err))
		return
	}

	if userModel.ComparePassword(validator.User.Password) != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("unauthorized", err))
		return
	}

	UpdateContextUserModel(c, userModel.ID)
	serializer := LoginSerializer{c}

	c.JSON(http.StatusOK, gin.H{"data": serializer.Response()})
}

func GetUser (c *gin.Context) {
	serializer := UserSerializer{c}

	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}