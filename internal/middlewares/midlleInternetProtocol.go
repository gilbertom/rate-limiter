package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

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
) *http.Request {
	ip := r.Context().Value("IP").(string)
	token := r.Context().Value("Token").(string)
	
	maxRequestsBySecond := getMaxRequestsBySecond(token, env.MaxRequestsAllowedByIP, env.MaxRequestsAllowedByToken)

	// Verifica se IP bloqueado
	key := ip + " - bloqueado"
	found, _, err := storageUseCase.FindByKey(key)
	if err != nil {
		log.Fatal("Error on find keys with prefix", err.Error())
		ct := context.WithValue(r.Context(), "isError", err)
		return r.WithContext(ct)
	}
	if found {
		ct := context.WithValue(r.Context(), "isBlock", true)
		return r.WithContext(ct)
	}

	// Increment the request count
	key = "rate_limiter - " + ip + " - " + token + strconv.Itoa(time.Now().Second())
	countRequest, err := storageUseCase.Incr(key)
	if err != nil {
		log.Fatal("Error on incremenet request", err.Error())
		ct := context.WithValue(r.Context(), "isError", err)
		return r.WithContext(ct)
	}

	// Set the expiration time for the key
	err = storageUseCase.Expire(key, 1)
	if err != nil {
		log.Fatal("Error on expire", err)
		ct := context.WithValue(r.Context(), "isError", err)
		return r.WithContext(ct)
	}

	if int(countRequest) == maxRequestsBySecond + 1 {
		key := ip + " - bloqueado"
		_, err := storageUseCase.Incr(key)
		if err != nil {
			log.Fatal("Error on incremenet request", err.Error())
			ct := context.WithValue(r.Context(), "isError", err)
			return r.WithContext(ct)
		}
		err = storageUseCase.Expire(key, getTimeToReleaseRequest(env))
		if err != nil {
			log.Fatal("Error on expire", err)
			ct := context.WithValue(r.Context(), "isError", err)
			return r.WithContext(ct)
		}
		ct := context.WithValue(r.Context(), "isBlock", true)
		return r.WithContext(ct)
	}

	ct := context.WithValue(r.Context(), "isBlock", false)
	return r.WithContext(ct)
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
