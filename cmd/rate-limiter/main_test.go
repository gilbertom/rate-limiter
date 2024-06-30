package main

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gilbertom/go-rate-limiter/internal/config"
	"github.com/joho/godotenv"
)

// TestRateLimiterSemExcederLimitePorIp
// TestRateLimiterExcedendoLimitePorIp

// TestRateLimiterSemExcederLimitePorToken
// TestRateLimiterExcedendoLimitePorToken


func TestRateLimiterSemExcederLimitePorIp(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP
	iterations := 10
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
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP + 1
	log.Println("rateLimit", rateLimit)
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

		// if resp.StatusCode == http.StatusTooManyRequests {
		// 	// t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, i+1)
		// 	break
		// }
		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}
}

func TestRateLimiterSemExcederLimitePorIp(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP
	iterations := 10
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
	err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByIP + 1
	log.Println("rateLimit", rateLimit)
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

		// if resp.StatusCode == http.StatusTooManyRequests {
		// 	// t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, i+1)
		// 	break
		// }
		resp.Body.Close()
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}
}