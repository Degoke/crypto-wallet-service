package common

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewValidationError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)

	for _, e := range errs {
		if e.Param() != "" {
			res.Errors[e.Field()] = fmt.Sprintf("{%v: %v}", e.Tag(), e.Param())
		} else {
			res.Errors[e.Field()] = fmt.Sprintf("{key: %s}", e.Tag())
		}
	}

	return res
}

func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()

	return res
}

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.MustBindWith(obj, b)
}

func GenToken(id uuid.UUID) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"expiry": time.Now().Add(time.Hour * 24).Unix(),
	})
	secret := GetEnv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func StringTowei(str string) *big.Int {
	amountFloat, _ := strconv.ParseFloat(str, 64)
	
	wei := new(big.Float).Mul(big.NewFloat(amountFloat), big.NewFloat(1e18))
	weiInt := new(big.Int)
	wei.Int(weiInt)
	return weiInt
}

func WeiToString(wei *big.Int) string {
	weiFloat := new(big.Float).SetInt(wei)
	bnbFloat := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	bnbStr := fmt.Sprintf("%.18f", bnbFloat)
	return bnbStr
}