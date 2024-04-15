package info_receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KseniiaSalmina/Car-catalog/internal/config"
	"github.com/KseniiaSalmina/Car-catalog/internal/models"
)

type Receiver struct {
	client *http.Client
	url    string
}

func NewReceiver(cfg config.Receiver) *Receiver {
	return &Receiver{
		client: &http.Client{},
		url:    cfg.URL,
	}
}

func (r *Receiver) GetCarInfo(ctx context.Context, regNum string) (*models.Car, error) {
	byteInfo, err := r.getCarInfo(ctx, regNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get car info: %w", err)
	}

	var car models.Car
	if err := json.Unmarshal(byteInfo, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car info: %w", err)
	}

	return &car, nil
}

func (r *Receiver) getCarInfo(ctx context.Context, regNum string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.url+"/info?regNum="+regNum, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get car info: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get car info: %w", err)
	}

	bodyText, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to get car info: %w", err)
	}

	return bodyText, nil
}
