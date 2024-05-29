package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gilbertom/go-rate-limiter/internal/dto"
	"github.com/go-redis/redis/v8"
)

// RateLimiter is a middleware that limits the rate of incoming requests.
func RateLimiter(ctx context.Context, rdb *redis.Client, env dto.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value("Token")
		ip := r.Context().Value("IP")
		key := "rate_limiter - " + ip.(string) + " - " + token.(string)

		// Increment the request count
		countRequest, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			log.Fatal("Error on incremenet request", err.Error())
			ct := context.WithValue(r.Context(), "isError", err)
			r = r.WithContext(ct)
		}

		// Set TTL for the key if it's the first request
		if countRequest == 1 {
			err := rdb.Expire(ctx, key, time.Second).Err()
			if err != nil {
				log.Fatal("Error on expire", err)
				ct := context.WithValue(r.Context(), "isError", err)
				r = r.WithContext(ct)
			}
		}
		
		maxRequestsBySecond := getMaxRequestsBySecond(env.MaxRequestsAllowedByIP, env.MaxRequestsAllowedByToken)
		
		if int(countRequest) == (maxRequestsBySecond + 1) {
			ct := context.WithValue(r.Context(), "isBlock", true)
			r = r.WithContext(ct)

			err = rdb.Expire(ctx, key, time.Second * time.Duration(env.TimeToReleaseRequests)).Err()
			if err != nil {
				log.Fatal("Error on expire Time to Release", err)
				ct := context.WithValue(r.Context(), "isError", err)
				r = r.WithContext(ct)
			}
		}
		
		if int(countRequest) > maxRequestsBySecond {
			ct := context.WithValue(r.Context(), "isBlock", true)
			r = r.WithContext(ct)
		} else {
			ct := context.WithValue(r.Context(), "isBlock", false)
			r = r.WithContext(ct)
		}
		next.ServeHTTP(w, r)
	})
}

func getMaxRequestsBySecond(maxReqAllowedByIP, maxReqAllowedByToken int) (maxRequestsBySecond int) {
	if maxReqAllowedByIP >= maxReqAllowedByToken {
		return maxReqAllowedByIP
	}
	return maxReqAllowedByToken
}
