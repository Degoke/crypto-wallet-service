package user

import (
	"fmt"
	"net/http"

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/google/uuid"
)

func stripBearerPrefixFromToken(token string) (string, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:], nil
	}
	return token, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

func UpdateContextUserModel(c *gin.Context, userId uuid.UUID) {
	var userModel User

	if userId != uuid.Nil {
		db := common.GetDB()
		db.First(&userModel, userId)
	}

	c.Set("userId", userId)
	c.Set("user", userModel)
}

func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, uuid.Nil)
		token, err := request.ParseFromRequest(c.Request, AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			secret := common.GetEnv("JWT_SECRET")
			return []byte(secret), nil
		})

		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims)
			stringId := claims["id"].(string)
			userId, err := uuid.Parse(stringId)
			if err != nil {
				if auto401 {
					c.AbortWithError(http.StatusUnauthorized, err)
				}
				return
			}

			UpdateContextUserModel(c, userId)
		} else {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

	}
}