package user

import (

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/gin-gonic/gin"
)

type UserValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=255"`
	} `json:"user"`
	userModel User `json:"-"`
}

func (u *UserValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, u)

	if err != nil {
		return err
	}

	u.userModel.Email = u.User.Email
	u.userModel.Password = u.User.Password

	return nil
}

func NewUserValidator() UserValidator {
	userValidator := UserValidator{}
	return userValidator
}

type LoginValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=255"`
	} `json:"user"`
	userModel User `json:"-"`
}

func (l *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, l)

	if err != nil {
		return err
	}

	l.userModel.Email = l.User.Email
	l.userModel.Password = l.User.Password

	return nil
}

func NewLoginValidator() LoginValidator {
	loginValidator := LoginValidator{}
	return loginValidator
}