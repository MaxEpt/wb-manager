package wb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wb-manager/config"
	"wb-manager/internal/dto"
)

type Wb struct {
	cfg        *config.WbApiConfig
	httpClient http.Client
}

func New(c *config.WbApiConfig) *Wb {
	httpClient := http.Client{
		Timeout: 1 * time.Second,
	}
	return &Wb{cfg: c, httpClient: httpClient}
}

func (wb *Wb) UpdatePrices(newPrices dto.WbPriceUpdateRequest) error {
	reqBody, err := json.Marshal(newPrices)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/public/api/v1/prices", wb.cfg.ApiUrl), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", wb.cfg.Token))

	resp, err := wb.httpClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resp code is %d, resp body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
