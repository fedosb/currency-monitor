package currency_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
)

type Gateway struct {
	ApiURL string

	client *http.Client
}

const rubJson = "rub.json"

type rubRatesResponse struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

func New(cfg config.CurrencyApiConfig) *Gateway {
	return &Gateway{
		ApiURL: cfg.GetURL(),
		client: &http.Client{},
	}
}

func (g *Gateway) GetRubRates(ctx context.Context) (dto.RubRates, error) {
	url := fmt.Sprintf("%s/%s", g.ApiURL, rubJson)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dto.RubRates{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return dto.RubRates{}, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.RubRates{}, fmt.Errorf("got http code: %d", resp.StatusCode)
	}

	var response rubRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.RubRates{}, fmt.Errorf("decoding response: %w", err)
	}

	date, err := time.Parse(time.DateOnly, response.Date)
	if err != nil {
		return dto.RubRates{}, fmt.Errorf("parsing date: %w", err)
	}

	return dto.RubRates{
		Date:  date,
		Rates: response.Rub,
	}, nil
}
