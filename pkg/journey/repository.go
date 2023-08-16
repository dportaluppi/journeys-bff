package journey

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Repository interface {
	Create(ctx context.Context, j *Journey) (*Journey, error)
	UniqueID() string
}

type HTTPRepo struct {
	endpoint string
	token    string
}

func NewHTTPRepo(endpoint, token string) *HTTPRepo {
	return &HTTPRepo{endpoint: endpoint, token: token}
}

func (hr *HTTPRepo) Create(ctx context.Context, j *Journey) (*Journey, error) {
	journeyBytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hr.endpoint+j.AccountID+"/journeys", bytes.NewBuffer(journeyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hr.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create journey, status: " + resp.Status)
	}

	var journeyResponse Journey
	if err := json.NewDecoder(resp.Body).Decode(&journeyResponse); err != nil {
		return nil, err
	}

	return &journeyResponse, nil
}

func (hr *HTTPRepo) UniqueID() string {
	objectID := primitive.NewObjectID()
	return objectID.Hex()
}
