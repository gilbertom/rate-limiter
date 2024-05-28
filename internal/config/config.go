package config

import (
	"os"

	"strconv"

	"github.com/gilbertom/go-rate-limiter/internal/dto"
)

// LoadEnv loads the environment variables and returns a dto.Env struct.
func LoadEnv() dto.Env{
	maxAllowedByIP, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_ALLOWED_BY_IP"))
	if err != nil {
		panic(err)
	}

	maxAllowedByToken, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_ALLOWED_BY_TOKEN"))
	if err != nil {
		panic(err)
	}

	limitRequestsByIP, err := strconv.ParseBool(os.Getenv("LIMIT_REQUESTS_BY_IP"))
	if err != nil {
		panic(err)
	}

	timeToReleaseRequests, err := strconv.Atoi(os.Getenv("TIME_TO_RELEASE_REQUESTS"))
	if err != nil {
		panic(err)
	}

	return dto.Env{
		ServerPort:                 os.Getenv("SERVER_PORT"),
		RedisAddr:                  os.Getenv("REDIS_ADDR"),
		MaxRequestsAllowedByIP:     maxAllowedByIP,
		MaxRequestsAllowedByToken:  maxAllowedByToken,
		LimitRequestsByIP:          limitRequestsByIP,
		OnRequestsExceededBlockBy:  os.Getenv("ON_REQUEST_EXCEEDED_BLOCK_BY"),
		TimeToReleaseRequests:      timeToReleaseRequests,
	}
}