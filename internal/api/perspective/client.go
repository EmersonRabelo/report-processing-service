package perspective

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"time"

	"github.com/EmersonRabelo/report-processing-service/internal/dto/report"
	dto "github.com/EmersonRabelo/report-processing-service/internal/dto/report"
)

type PerspectiveAPIClient interface {
	AnalyzePost(comment *string) (dto.PerspectiveAPIResponse, error)
}

type perspectiveAPIClient struct {
	httpClient *http.Client
	baseURL    string
}

type Comment struct {
	Text string `json:"text"`
}

type CommentPayload struct {
	Comment             Comment             `json:"comment"`
	RequestedAttributes map[string]struct{} `json:"requestedAttributes"`
}

func NewPerspectiveAPIClient(baseURL string) PerspectiveAPIClient {
	return &perspectiveAPIClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

func (client *perspectiveAPIClient) AnalyzePost(body *string) (dto.PerspectiveAPIResponse, error) {
	type Comment struct {
		Text string `json:"text"`
	}
	type CommentPayload struct {
		Comment             Comment             `json:"comment"`
		RequestedAttributes map[string]struct{} `json:"requestedAttributes"`
	}

	comment := Comment{
		Text: *body,
	}
	requestedAttributes := map[string]struct{}{
		"TOXICITY": {},
	}
	data := CommentPayload{
		Comment:             comment,
		RequestedAttributes: requestedAttributes,
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return dto.PerspectiveAPIResponse{}, err
	}
	req, err := http.NewRequest("POST", client.baseURL, bytes.NewBuffer(buf))
	if err != nil {
		return dto.PerspectiveAPIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return dto.PerspectiveAPIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return dto.PerspectiveAPIResponse{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(b))
	}

	var apiResp report.PerspectiveAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return report.PerspectiveAPIResponse{}, err
	}

	return apiResp, nil
}
