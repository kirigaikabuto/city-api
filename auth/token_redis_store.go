package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"os"
	"time"
)

const (
	accessKeyName  = "access_key"
	refreshKeyName = "refresh_key"
	accessSecret   = "jdnfksdmfksd"
	refreshSecret  = "mcmvmkmsdnfsdmfdsjf"
)

type tokenStore struct {
	redisClient *redis.Client
}

func NewTokenStore(config RedisConfig) (TokenStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		DB:       0,
		Password: config.Password,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &tokenStore{redisClient: client}, nil
}

func (t *tokenStore) CreateToken(cmd *CreateTokenCommand) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 45).Unix()
	td.AccessUuid = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()
	var err error
	_ = os.Setenv("ACCESS_SECRET", accessSecret)
	_ = os.Setenv("REFRESH_SECRET", refreshSecret)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = cmd.UserId
	atClaims["exp"] = td.AtExpires
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_type"] = cmd.UserType
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = cmd.UserId
	rtClaims["exp"] = td.RtExpires
	rtClaims["user_type"] = cmd.UserType
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	//save token
	aT := time.Unix(td.AtExpires, 0)
	rT := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err = t.redisClient.Set(accessKeyName+":"+td.AccessUuid, cmd.UserId, aT.Sub(now)).Err()
	if err != nil {
		return nil, err
	}
	err = t.redisClient.Set(refreshKeyName+":"+td.RefreshUuid, cmd.UserId, rT.Sub(now)).Err()
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (t *tokenStore) GetToken(id string) (string, error) {
	userId, err := t.redisClient.Get(accessKeyName + ":" + id).Result()
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (t *tokenStore) RemoveToken(id string) (int64, error) {
	deleted, err := t.redisClient.Del(accessKeyName + ":" + id).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
