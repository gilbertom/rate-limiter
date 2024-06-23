package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gilbertom/go-rate-limiter/internal/dto"
	usecases "github.com/gilbertom/go-rate-limiter/internal/usecase"
)

// RateLimiter is a middleware that limits the rate of incoming requests.
func RateLimiter(
	ctx context.Context, 
	storageUseCase *usecases.StorageUseCase, 
	env dto.Env, 
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = processHateLimiter(r, env, storageUseCase)
		next.ServeHTTP(w, r)
	})
}

func processHateLimiter(
	r *http.Request,
	env dto.Env,
	storageUseCase *usecases.StorageUseCase,
) *http.Request{
	ip := r.Context().Value("IP").(string)
	token := r.Context().Value("Token").(string)
	key := "rate_limiter - " + ip + " - " + token
	maxRequestsBySecond := getMaxRequestsBySecond(token, env.MaxRequestsAllowedByIP, env.MaxRequestsAllowedByToken)

	// Increment the request count
	countRequest, err := storageUseCase.Incr(key)
	if err != nil {
		log.Fatal("Error on incremenet request", err.Error())
		ct := context.WithValue(r.Context(), "isError", err)
		r = r.WithContext(ct)
	}

	// Set TTL for the key if it's the first request
	if countRequest == 1 {
		err := storageUseCase.Expire(key, maxRequestsBySecond)
		if err != nil {
			log.Fatal("Error on expire", err)
			ct := context.WithValue(r.Context(), "isError", err)
			r = r.WithContext(ct)
		}
	}
	
	if int(countRequest) == (maxRequestsBySecond + 1) {
		ct := context.WithValue(r.Context(), "isBlock", true)
		r = r.WithContext(ct)

		err = storageUseCase.Expire(key, getTimeToReleaseRequest(env))
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
	return r
}

func getMaxRequestsBySecond(token string, maxReqAllowedByIP, maxReqAllowedByToken int) (maxRequestsBySecond int) {
	if token != "token not present" && maxReqAllowedByToken >= maxReqAllowedByIP {
		return maxReqAllowedByToken
	}
	return maxReqAllowedByIP
}

func getTimeToReleaseRequest(env dto.Env) (timeToReleaseRequests int) {
	if env.OnRequestsExceededBlockBy == "IP" {
		return env.TimeToReleaseRequestsIP
	}
	return env.TimeToReleaseRequestsToken
}
