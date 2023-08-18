package journey

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
)

type Repository interface {
	Create(ctx context.Context, j *JourneyWriteModel) (*JourneyWriteModel, error)
	GetJourneys(ctx context.Context, accountId string) ([]JourneyWriteModel, error)
	GetJourneyByID(ctx context.Context, accountId, journeyID string) (*JourneyWriteModel, error)
	UniqueID() string
}

type HTTPRepo struct {
	endpoint string
	token    string
}

func NewHTTPRepo(endpoint, token string) *HTTPRepo {
	return &HTTPRepo{endpoint: endpoint, token: token}
}

func (hr *HTTPRepo) Create(ctx context.Context, j *JourneyWriteModel) (*JourneyWriteModel, error) {
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create journey. Status: " + resp.Status + ". Body: " + string(body))
	}

	var journeyResponse JourneyWriteModel
	if err := json.Unmarshal(body, &journeyResponse); err != nil {
		return nil, err
	}

	return &journeyResponse, nil
}

func (hr *HTTPRepo) GetJourneys(ctx context.Context, accountId string) ([]JourneyWriteModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", hr.endpoint+accountId+"/journeys", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+hr.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch journeys. Status: " + resp.Status + ". Body: " + string(body))
	}

	var journeys []JourneyWriteModel
	if err := json.Unmarshal(body, &journeys); err != nil {
		return nil, err
	}

	return journeys, nil
}

func (hr *HTTPRepo) GetJourneyByID(ctx context.Context, accountId, journeyID string) (*JourneyWriteModel, error) {
	// Build the full URL to fetch the journey by its ID
	url := hr.endpoint + accountId + "/journeys/" + journeyID

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+hr.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, _ := io.ReadAll(resp.Body)

	// Check the HTTP status code to determine if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch journey by ID. Status: " + resp.Status + ". Body: " + string(body))
	}

	var journey JourneyWriteModel
	if err := json.Unmarshal(body, &journey); err != nil {
		return nil, err
	}

	return &journey, nil
}

func (hr *HTTPRepo) UniqueID() string {
	objectID := primitive.NewObjectID()
	return objectID.Hex()
}
