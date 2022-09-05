package mdw

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"strings"
)

type Middleware struct {
	tokenStore TokenStore
}

func NewMiddleware(tkn TokenStore) Middleware {
	return Middleware{tokenStore: tkn}
}

func (m *Middleware) MakeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = os.Setenv("ACCESS_SECRET", accessSecret)
		_ = os.Setenv("REFRESH_SECRET", refreshSecret)
		tokenAuth, err := m.ExtractTokenMetadata(c.Request)
		if err != nil && err == redis.Nil {
			respondJSONMdw(c.Writer, http.StatusBadRequest, "your token is expired")
			c.Abort()
			return
		} else if err != nil {
			respondJSONMdw(c.Writer, http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", tokenAuth.UserId)
		c.Set("access_uuid", tokenAuth.AccessUuid)
		c.Set("user_type", tokenAuth.UserType)
		c.Next()
	}
}

func (m *Middleware) verifyToken(r *http.Request) (*jwt.Token, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("need bearer token")
	}
	tokenString := strArr[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Middleware) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := m.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("not access uuid")
		}
		userId := claims["user_id"].(string)
		userType := claims["user_type"].(string)
		_, err = m.tokenStore.GetToken(accessUuid)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
			UserType:   userType,
		}, nil
	}
	return nil, errors.New("error during extract of token metadata")
}

func respondJSONMdw(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
