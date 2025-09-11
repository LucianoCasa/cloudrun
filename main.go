package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	cepBaseURL     = "https://viacep.com.br/ws"
	weatherBaseURL = "http://api.weatherapi.com/v1"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

func main() {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	cepSvc := &CepService{BaseURL: cepBaseURL, HTTPClient: httpClient}
	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")
	weatherSvc := &WeatherService{BaseURL: weatherBaseURL, APIKey: weatherAPIKey, HTTPClient: httpClient}

	http.HandleFunc("/", weatherHandler(cepSvc, weatherSvc))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func weatherHandler(cepSvc *CepService, weatherSvc *WeatherService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		if !cepRegex.MatchString(cep) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}
		city, err := cepSvc.Lookup(r.Context(), cep)
		if err != nil {
			if errors.Is(err, ErrCepNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("can not find zipcode"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tempC, err := weatherSvc.GetTempC(r.Context(), city)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := map[string]float64{
			"temp_C": tempC,
			"temp_F": tempC*1.8 + 32,
			"temp_K": tempC + 273,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// CepService calls ViaCEP to get city name from cep

type CepService struct {
	BaseURL    string
	HTTPClient *http.Client
}

var ErrCepNotFound = errors.New("cep not found")

func (c *CepService) Lookup(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("%s/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep status: %d", resp.StatusCode)
	}
	var data struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.Erro || data.Localidade == "" {
		return "", ErrCepNotFound
	}
	return data.Localidade, nil
}

// WeatherService calls WeatherAPI to get temperature in Celsius

type WeatherService struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func (wSvc *WeatherService) GetTempC(ctx context.Context, city string) (float64, error) {
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", wSvc.BaseURL, wSvc.APIKey, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := wSvc.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status: %d", resp.StatusCode)
	}
	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}
