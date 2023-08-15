package audience

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPAudienceRepo struct {
	baseURL string
	token   string
}

func NewHTTPAudienceRepo(baseURL, token string) *HTTPAudienceRepo {
	return &HTTPAudienceRepo{baseURL: baseURL, token: token}
}

func (repo *HTTPAudienceRepo) GetAudiences(ctx context.Context, filter *Filter, pageSize int, pageNumber int) (*AudiencesResponse, error) {
	var repoResponse RepoAudiencesResponse

	url := repo.baseURL + fmt.Sprintf("?provider=%s&page=%d&pageSize=%d", filter.Provider, pageNumber, pageSize)
	if filter.Query != "" {
		url += "&name=" + filter.Query
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+repo.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &repoResponse)
	if err != nil {
		return nil, err
	}

	domainAudiences := make([]Audience, len(repoResponse.Data))
	totalCustomers := 0
	for i, repoAud := range repoResponse.Data {
		domainAudiences[i] = toDomain(repoAud)
		totalCustomers += domainAudiences[i].LastGeneratedSize
	}

	response := &AudiencesResponse{
		Data: Data{
			Audiences: domainAudiences,
			Customers: CustomersInfo{Total: totalCustomers},
		},
		Metadata: repoResponse.Metadata,
	}

	return response, nil
}

func toDomain(repoAud RepoAudience) Audience {
	return Audience{
		ID:   repoAud.ID,
		Name: repoAud.Name,
		Filters: AudienceFilters{
			CustomerAttributesFilter: fmt.Sprintf("%v", repoAud.Filters.CustomerAttributesFilter),
			EventsFilter:             fmt.Sprintf("%v", repoAud.Filters.EventsFilter),
		},
		Tags:              repoAud.Tags,
		LastGeneratedSize: repoAud.LastGeneratedSize,
		UpdatedAt:         repoAud.UpdatedAt,
	}
}

// --------- Repository model -------------

type RepoAudience struct {
	ID                string      `json:"id"`
	Provider          string      `json:"provider"`
	Name              string      `json:"name"`
	Filters           RepoFilters `json:"filters"`
	Tags              []string    `json:"tags"`
	LastGeneratedSize int         `json:"last_generated_size"`
	UpdatedAt         string      `json:"last_generated_at"`
	Folder            string      `json:"folder"`
}

type RepoFilters struct {
	CustomerAttributesFilter map[string]interface{} `json:"customerAttributesFilter"`
	EventsFilter             map[string]interface{} `json:"eventsFilter"`
}

type RepoAudiencesResponse struct {
	Data     []RepoAudience `json:"data"`
	Metadata Metadata       `json:"metadata"`
}
