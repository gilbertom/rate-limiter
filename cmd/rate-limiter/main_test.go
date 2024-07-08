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

func TestRateLimiterSemExcederLimitePorIp(t *testing.T) {
	log.Println("Iniciando teste de rate limiter sem exceder limite por IP")
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
	log.Println("Fim do teste de rate limiter sem exceder limite por IP")
}


func TestRateLimiterExcedendoLimitePorIp(t *testing.T) {
	log.Println("Iniciando teste de rate limiter excedendo limite por IP")
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
		log.Println("Excedido limite de requests por IP")
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}
	log.Println("Fim do teste de rate limiter excedendo limite por IP")
}

func TestRateLimiterSemExcederLimiteUsandoToken(t *testing.T) {
	log.Println("Iniciando teste de rate limiter sem exceder limite por Token")
	err := godotenv.Load("/cmd/rate-limiter/.env")
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
	log.Println("Fim do teste de rate limiter sem exceder limite por Token")
}

func TestRateLimiterExcedendoLimiteUsandoToken(t *testing.T) {
	log.Println("Iniciando teste de rate limiter excedendo limite por Token")
	err := godotenv.Load("/cmd/rate-limiter/.env")
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
		log.Println("Excedido limite de requests por Token")
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, 1)
	}
	log.Println("Fim do teste de rate limiter excedendo limite por Token")
}

func TestRateLimiterExcedendoLimitePorIpAguardaLiberacaoBloqueio(t *testing.T) {
	log.Println("Iniciando teste de rate limiter excedendo limite por IP e aguardando liberação do bloqueio")
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

	timeWait := env.TimeToReleaseRequestsIP + 0
	log.Printf("Aguardando %d segundos até ocorrer o desbloqueio por IP", timeWait)
	time.Sleep(time.Duration(timeWait)* time.Second)
	log.Printf("Aguardado %d segundos", timeWait)

	resp, err = client.Get(serverURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v on request %d", resp.StatusCode, rateLimit)
	}
	log.Println("Fim do teste de rate limiter excedendo limite por IP e aguardando liberação do bloqueio")
}

func TestRateLimiterExcedendoLimitePorTokenAguardaLiberacaoBloqueio(t *testing.T) {
	log.Println("Iniciando teste de rate limiter excedendo limite por Token e aguardando liberação do bloqueio")
	err := godotenv.Load("/cmd/rate-limiter/.env")
	// err := godotenv.Load()
	if err != nil {
			log.Fatalf("Error trying to load env variables: %v", err)
	}
	env := config.LoadEnv()
		
	rateLimit := env.MaxRequestsAllowedByToken + 1
	serverURL := "http://localhost:8080/"

	os.Setenv("ON_REQUESTS_EXCEEDED_BLOCK_BY", "TOKEN")
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
		t.Errorf("Expected status 429 but got %v on request %d", resp.StatusCode, rateLimit)
	}

	timeWait := env.TimeToReleaseRequestsToken + 0
	log.Printf("Aguardando %d segundos até ocorrer o desbloqueio por Token", timeWait)
	time.Sleep(time.Duration(timeWait)* time.Second)
	log.Printf("Aguardado %d segundos", timeWait)

	resp, err = client.Get(serverURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v on request %d", resp.StatusCode, rateLimit)
	}
	log.Println("Fim do teste de rate limiter excedendo limite por Token e aguardando liberação do bloqueio")
}