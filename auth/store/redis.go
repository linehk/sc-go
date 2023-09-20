package store

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/stablecog/sc-go/database"
)

type AuthorizationRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
}

type RedisStore struct {
	Ctx         context.Context
	RedisClient *database.RedisWrapper
}

func NewRedisStore(ctx context.Context) *RedisStore {
	// Setup redis
	redis, err := database.NewRedis(ctx)
	if err != nil {
		log.Fatal("Error connecting to redis", "err", err)
		os.Exit(1)
	}

	return &RedisStore{
		Ctx:         ctx,
		RedisClient: redis,
	}
}

func (s *RedisStore) SaveAuthRequestInCache(authReq *AuthorizationRequest) error {
	// Serialize auth req
	marshalled, err := json.Marshal(authReq)
	if err != nil {
		return err
	}
	return s.RedisClient.Client.Set(s.Ctx, authReq.Code, marshalled, 5*time.Minute).Err()
}

func (s *RedisStore) GetAuthRequestFromCache(code string) (*AuthorizationRequest, error) {
	// Get auth req from cache
	val, err := s.RedisClient.Client.Get(s.Ctx, code).Result()
	if err != nil {
		return nil, err
	}
	var authReq AuthorizationRequest
	err = json.Unmarshal([]byte(val), &authReq)
	if err != nil {
		return nil, err
	}
	return &authReq, nil
}
