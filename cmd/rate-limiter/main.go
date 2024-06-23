package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gilbertom/go-rate-limiter/internal/config"
	formatresponse "github.com/gilbertom/go-rate-limiter/internal/formatResponse"
	"github.com/gilbertom/go-rate-limiter/internal/infra"
	"github.com/gilbertom/go-rate-limiter/internal/middlewares"
	usecases "github.com/gilbertom/go-rate-limiter/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}

	env := config.LoadEnv()

	redisStorage, err := infra.NewRedisStorage(env.RedisAddr)
	if err != nil {
			log.Fatalf("Error while creating Redis storage: %v", err)
	}

	redisUseCase := usecases.NewStorageUseCase(redisStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ResponseHandler(ctx, w, r)
	})

	wrappedHandler := middlewares.GetTokenFromHeader(
    middlewares.GetIPfromClient(
        middlewares.ProcessHandler(
            middlewares.RateLimiter(ctx, redisUseCase, env, mux),
        ),
    ),
	)

	fmt.Println("Listening on port", env.ServerPort)
	http.ListenAndServe(":" + env.ServerPort, wrappedHandler)
}

// ResponseHandler handles the HTTP response to the client.
func ResponseHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	msg := formatresponse.ResponseOutput{}

	isError := r.Context().Value("isError")
	isBlock := r.Context().Value("isBlock")

	if err, ok := isError.(bool); ok && err {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	if err, ok := isBlock.(bool); ok && err {
		w.WriteHeader(http.StatusTooManyRequests)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")

	// Get the response
	responseStruct := msg.GetResponse(isBlock.(bool))

	jsonResponse, err := json.Marshal(responseStruct)
	if err != nil {
			http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
			return
	}
	w.Write(jsonResponse)
}
