package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

var _ ModelManager = &OllamaManager{}

type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done,omitempty"`
	DoneReason         string `json:"done_reason,omitempty"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int    `json:"total_duration,omitempty"`
	LoadDuration       int    `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int    `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int    `json:"eval_duration,omitempty"`
}

type EmbedResponse struct {
	Model           string        `json:"model"`
	Embeddings      [][]float32   `json:"embeddings"`
	TotalDuration   time.Duration `json:"total_duration,omitempty"`
	LoadDuration    time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount int           `json:"prompt_eval_count,omitempty"`
}

type OllamaManager struct {
	model   string
	baseUrl string
}

// Embed generates embeddings from a model.
func (o *OllamaManager) Embed(ctx context.Context, datas []string) ([][]float32, error) {
	log := logrus.WithContext(ctx)
	client := resty.New()
	body := map[string]interface{}{
		"model": o.model,
		"input": datas,
	}
	url := fmt.Sprintf("%s%s", o.baseUrl, "/api/embed")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		log.Errorf("post %s failed, err: %v", url, err)
		return nil, nil
	}
	embedResp := &EmbedResponse{}
	if err := json.Unmarshal(resp.Body(), &embedResp); err != nil {
		log.Errorf("json.Unmarshal EmbedResponse failed, err: %v", err)
		return nil, err
	}
	return embedResp.Embeddings, nil
}

// Generate generates a response for a given prompt.
func (o *OllamaManager) Generate(ctx context.Context, prompt string) (string, error) {
	log := logrus.WithContext(ctx)
	client := resty.New()
	body := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
	}
	url := fmt.Sprintf("%s%s", o.baseUrl, "/api/generate")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		log.Errorf("post %s failed, err: %v", url, err)
		return "", nil
	}
	genResp := &GenerateResponse{}
	if err := json.Unmarshal(resp.Body(), &genResp); err != nil {
		log.Errorf("json.Unmarshal GenerateResponse failed, err: %v", err)
		return "", err
	}
	return genResp.Response, nil
}

func NewOllamaManager(model, url string) *OllamaManager {
	return &OllamaManager{
		model:   model,
		baseUrl: url,
	}
}
