package user

import (
	"time"

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserSerializer struct {
	C *gin.Context
}

type LoginSerializer struct {
	C *gin.Context
}

type UserResponse struct {
	ID    uuid.UUID   `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Duration `json:"created_at"`
	UpdatedAt time.Duration `json:"updated_at"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (u *UserSerializer) Response() UserResponse {
	user := u.C.MustGet("user").(User)
	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		CreatedAt: time.Since(user.CreatedAt),
		UpdatedAt: time.Since(user.UpdatedAt),
	}
	return response
}

func (u *LoginSerializer) Response() loginResponse {
	user := u.C.MustGet("user").(User)
	response := loginResponse{
		Token: common.GenToken(user.ID),
	}
	return response
}