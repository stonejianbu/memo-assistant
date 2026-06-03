package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stonejianbu/memo-assistant/pkg/llm"
)

type DouBaoEmbedResponse struct {
	Data *DouBaoEmbedData `json:"data"`
}

type DouBaoEmbedData struct {
	Embedding []float32 `json:"embedding"`
}

type DouBaoResponse struct {
	ID              string             `json:"id"`
	Object          string             `json:"object"`
	Model           string             `json:"model"`
	Status          string             `json:"status"`
	ServiceTier     string             `json:"service_tier"`
	CreatedAt       int64              `json:"created_at"`
	ExpireAt        int64              `json:"expire_at"`
	MaxOutputTokens int                `json:"max_output_tokens"`
	Store           bool               `json:"store"`
	Output          []DouBaoOutputItem `json:"output"`
	Usage           DouBaoUsage        `json:"usage"`
	Caching         DouBaoCaching      `json:"caching"`
}

type DouBaoOutputItem struct {
	ID      string           `json:"id"`
	Type    string           `json:"type"`
	Status  string           `json:"status"`
	Role    string           `json:"role,omitempty"`
	Summary []DouBaoTextItem `json:"summary,omitempty"`
	Content []DouBaoTextItem `json:"content,omitempty"`
}

type DouBaoTextItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type DouBaoUsage struct {
	InputTokens         int                      `json:"input_tokens"`
	OutputTokens        int                      `json:"output_tokens"`
	TotalTokens         int                      `json:"total_tokens"`
	InputTokensDetails  DouBaoInputTokenDetails  `json:"input_tokens_details"`
	OutputTokensDetails DouBaoOutputTokenDetails `json:"output_tokens_details"`
}

type DouBaoInputTokenDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type DouBaoOutputTokenDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type DouBaoCaching struct {
	Type string `json:"type"`
}

type DouBao struct {
	apiKey string
}

// Embed generates embeddings from a model.
func (that *DouBao) Embed(ctx context.Context, datas []map[string]interface{}) ([]float32, error) {
	log := logrus.WithContext(ctx)
	embedResp := &DouBaoEmbedResponse{}
	client := NewDouBaoClient(
		"https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal",
		"doubao-embedding-vision-250615",
		that.apiKey,
	)
	if err := client.Post(ctx, map[string]interface{}{
		"input": datas,
	}, embedResp); err != nil {
		log.Errorf("client.Post failed, err: %v", err)
		return nil, err
	}
	if embedResp != nil && embedResp.Data != nil {
		return embedResp.Data.Embedding, nil
	}
	log.Warnf("embedResp: %+v", embedResp)
	return nil, nil
}

// Generate generates a response for a given prompt.
func (that *DouBao) Generate(ctx context.Context, prompt string) (string, error) {
	log := logrus.WithContext(ctx)
	resp := &DouBaoResponse{}
	client := NewDouBaoClient(
		"https://ark.cn-beijing.volces.com/api/v3/responses",
		"doubao-seed-2-0-pro-260215",
		that.apiKey,
	)
	if err := client.Post(ctx, map[string]interface{}{
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "input_text",
						"text": prompt,
					},
				},
			},
		},
	}, resp); err != nil {
		log.Errorf("client.Post failed, err: %v", err)
		return "", err
	}
	for _, item := range resp.Output {
		if item.Type == "message" {
			for _, c := range item.Content {
				if c.Type == "output_text" {
					return c.Text, nil
				}
			}
		}
	}
	return "", errors.New("no output_text found in response")
}

func NewDouBao(apiKey string) llm.ModelManager {
	return &DouBao{apiKey: apiKey}
}

type DouBaoClient struct {
	url    string
	model  string
	apiKey string
}

func NewDouBaoClient(url, model, apiKey string) *DouBaoClient {
	return &DouBaoClient{
		url:    url,
		model:  model,
		apiKey: apiKey,
	}
}

func (that *DouBaoClient) Post(ctx context.Context, reqData map[string]interface{}, result interface{}) error {
	log := logrus.WithContext(ctx)
	if reqData == nil {
		return errors.New("reqData is nil")
	}
	reqData["model"] = that.model
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", that.apiKey)).
		SetBody(reqData).
		Post(that.url)
	if err != nil {
		log.Errorf("post %s failed, err: %v", that.url, err)
		return err
	}
	if err := json.Unmarshal(resp.Body(), result); err != nil {
		log.Errorf("json.Unmarshal failed, err: %v", err)
		return err
	}
	return nil
}
