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
	"github.com/joho/godotenv"
)

func main() {
	var ctx = context.Background()
	
	// Loads environment variables from the file . env
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error loading file . env: %v", err)
	}

	env := config.LoadEnv()

	// Create a new Redis Client
	rdb := infra.NewRedisClient(env.RedisAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ResponseHandler(ctx, w, r)
	})

	wrappedHandler := middlewares.GetTokenFromHeader(
    middlewares.GetIPfromClient(
        middlewares.ProcessHandler(
            middlewares.RateLimiter(ctx, rdb, env, mux),
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
