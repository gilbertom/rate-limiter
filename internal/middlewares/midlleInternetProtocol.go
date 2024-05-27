package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiter is a middleware that limits the rate of incoming requests.
func RateLimiter(ctx context.Context, rdb *redis.Client, maxAllowedByIP int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Context().Value("IP")
		if ip == nil {
			fmt.Println("IP == nil")
		}
		
		key := "rate_limiter:" + ip.(string)

		// Increment the request count
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			log.Fatal("Erro on incremenet", err)
			ct := context.WithValue(r.Context(), "isError", err)
			r = r.WithContext(ct)
		}

		// Set TTL for the key if it's the first request
		if count == 1 {
			err := rdb.Expire(ctx, key, time.Second).Err()
			if err != nil {
				log.Fatal("Error on expire", err)
				ct := context.WithValue(r.Context(), "isError", err)
				r = r.WithContext(ct)
			}
		}
		
		if count > maxAllowedByIP {
			ct := context.WithValue(r.Context(), "isBlock", true)
			r = r.WithContext(ct)
		} else {
			ct := context.WithValue(r.Context(), "isBlock", false)
			r = r.WithContext(ct)
		}
		next.ServeHTTP(w, r)
	})
}
