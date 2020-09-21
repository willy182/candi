package auth

import (
	"context"

	"agungdwiprasetyo.com/backend-microservices/pkg/shared"
	"github.com/gomodule/redigo/redis"
)

type authServiceRedis struct {
	redisPool *redis.Pool
}

// NewAuthServiceRedis using redis
func NewAuthServiceRedis(redisHost string, redisPort int, redisAuth string, useTLS bool) Auth {
	return &authServiceRedis{}
}

func (a *authServiceRedis) ValidateToken(ctx context.Context, token string) (cl *shared.TokenClaim, err error) {

	return &shared.TokenClaim{}, nil
}

func (a *authServiceRedis) GenerateToken(ctx context.Context, claim *shared.TokenClaim) <-chan shared.Result {
	output := make(chan shared.Result)

	go func() { defer close(output) }()

	return output
}
