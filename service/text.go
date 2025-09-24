package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/dao"
	"github.com/stonejianbu/memo-assistant/pkg/llm"
	"strings"
)

var SystemPromptTemplate = `
Please answer user's question according to context,
If the Context is not empty, simply return the answer directly without excessive explanation.
Question: 
%s

Context:
%s
`

type DataManager interface {
	// Train data to weaviate
	Train(ctx context.Context, datas []string) error
	// Query approximation data from weaviate, generate a whole prompt send to llm generate answer
	Query(ctx context.Context, query string) (string, error)
}

type TextManager struct {
	llmClient llm.ModelManager
	class     string
}

// Train data to weaviate
func (d *TextManager) Train(ctx context.Context, datas []string) error {
	log := logrus.WithContext(ctx)
	log.Infof("Train, datas length: %d", len(datas))
	embedding, err := d.llmClient.Embed(ctx, datas)
	if err != nil {
		return err
	}
	for i, data := range datas {
		go func(index int, data string) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic: %v", err)
				}
			}()
			obj := map[string]interface{}{
				"content": data,
			}
			// save data to weaviate
			if err := dao.Create(ctx, d.class, obj, embedding[index]); err != nil {
				log.Errorf("dao.Create failed, err: %v", err)
			}
		}(i, data)
	}
	return nil
}

// Query approximation data from weaviate, generate a whole prompt send to llm generate answer
func (d *TextManager) Query(ctx context.Context, prompt string) (string, error) {
	log := logrus.WithContext(ctx)
	log.Infof("Query, Prompt: %s", prompt)
	log.Infof("generate embedding, prompt: %s", prompt)
	Embeddings, err := d.llmClient.Embed(ctx, []string{prompt})
	if err != nil {
		log.Errorf("llmClient.Embed failed, err: %v", err)
		return "", err
	}
	log.Infof("query approximation data from weaviate")
	results, err := dao.Query(ctx, d.class, prompt, Embeddings[0])
	if err != nil {
		log.Errorf("dao.Query failed, err: %v", err)
		return "", err
	}
	// generate a whole prompt with context.
	newPrompt := fmt.Sprintf(SystemPromptTemplate, prompt, strings.Join(results, "\n"))
	log.Infof("use llm generate text to answer the question, prompt:%s", newPrompt)
	answer, err := d.llmClient.Generate(ctx, newPrompt)
	if err != nil {
		log.Errorf("llmClient.Generate failed, err: %v", err)
		return "", err
	}
	return answer, nil
}

func NewTextManager(llmClient llm.ModelManager) *TextManager {
	return &TextManager{
		llmClient: llmClient,
		class:     config.Cfg.Weaviate.Class,
	}
}
