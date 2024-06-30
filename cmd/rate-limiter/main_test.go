package main

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gilbertom/go-rate-limiter/internal/config"
	"github.com/joho/godotenv"
)

// TestRateLimiterSemExcederLimitePorIp - feito
// TestRateLimiterExcedendoLimitePorIp - feito

// TestRateLimiterSemExcederLimitePorToken
// TestRateLimiterExcedendoLimitePorToken

// TestRateLimiterTokenSobrepondoLimitePorIP

// TestRateLimiterExcedendoLimitePorIpAguardaLiberacaoBloqueio
// TestRateLimiterExcedendoLimitePorTokenAguardaLiberacaoBloqueio

func TestRateLimiterSemExcederLimitePorIp(t *testing.T) {
	err := godotenv.Load("/cmd/rate-limiter/.env")
	if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
	}
		
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP
	iterations := 3
	serverURL := "http://localhost:8080/"

	go func() {
		main()
	}()

	time.Sleep(time.Second)

	client := &http.Client{}

	loopMain:
	for i := 1; i <= iterations; i++ {

		for i := 1; i <= rateLimit; i++ {
			resp, err := client.Get(serverURL)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
				break
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK but got %v on request %d", resp.StatusCode, i+1)
				break loopMain
			}
			resp.Body.Close()
		}
		time.Sleep(1*time.Second)
	}
}


func TestRateLimiterExcedendoLimitePorIp(t *testing.T) {
	err := godotenv.Load("/cmd/rate-limiter/.env")
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP + 1
	serverURL := "http://localhost:8080/"

	go func() {
		main()
	}()

	time.Sleep(time.Second)

	client := &http.Client{}

	var resp *http.Response

	for i := 1; i <= rateLimit; i++ {
		resp, err = client.Get(serverURL)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
			break
		}

		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}
}

func TestRateLimiterSemExcederLimiteUsandoToken(t *testing.T) {
	err := godotenv.Load("/cmd/rate-limiter/.env")
	// err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	
	os.Setenv("ON_REQUESTS_EXCEEDED_BLOCK_BY", "TOKEN")
	env := config.LoadEnv()
	
	rateLimit := env.MaxRequestsAllowedByToken
	iterations := 2
	serverURL := "http://localhost:8080/"
	apiKey := "123456"

	go func() {
		main()
	}()

	time.Sleep(time.Second)

	client := &http.Client{}

loopMain:
	for i := 1; i <= iterations; i++ {
		for i := 1; i <= rateLimit; i++ {
			req, err := http.NewRequest("GET", serverURL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("API_KEY", apiKey)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
				break
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK but got %v on request %d", resp.StatusCode, i+1)
				break loopMain
			}
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
}

func TestRateLimiterExcedendoLimiteUsandoToken(t *testing.T) {
	err := godotenv.Load("/cmd/rate-limiter/.env")
	// err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	
	os.Setenv("ON_REQUESTS_EXCEEDED_BLOCK_BY", "TOKEN")
	env := config.LoadEnv()
	
	rateLimit := env.MaxRequestsAllowedByToken + 1
	serverURL := "http://localhost:8080/"
	apiKey := "123456"

	go func() {
		main()
	}()

	time.Sleep(time.Second)

	client := &http.Client{}
	var resp *http.Response

	for i := 1; i <= rateLimit; i++ {
		req, err := http.NewRequest("GET", serverURL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("API_KEY", apiKey)
		
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
			break
		}
		defer resp.Body.Close()
	}
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, 1)
	}
}

func TestRateLimiterExcedendoLimitePorIpAguardaLiberacaoBloqueio(t *testing.T) {
	err := godotenv.Load("/cmd/rate-limiter/.env")
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP + 1
	serverURL := "http://localhost:8080/"

	go func() {
		main()
	}()

	time.Sleep(time.Second)

	client := &http.Client{}

	var resp *http.Response

	for i := 1; i <= rateLimit; i++ {
		resp, err = client.Get(serverURL)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
			break
		}

		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}

	log.Printf("Aguardando %d segundos até ocorrer o desbloqueio por IP", env.TimeToReleaseRequestsIP)
	time.Sleep(time.Duration(env.TimeToReleaseRequestsIP))
}