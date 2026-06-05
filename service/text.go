package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/dao"
	"github.com/stonejianbu/memo-assistant/pkg/llm"
	"strings"
	"sync"
)

var SystemPromptTemplate = `
请回答用户问题，优先基于Context内容来回答，如果Context为空或者无匹配，则基于事实和你已知的知识来回答，回答结果以markdown格式返回。
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
	wg := sync.WaitGroup{}
	for i, data := range datas {
		wg.Add(1)
		go func(index int, data string) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic: %v", err)
				}
				wg.Done()
			}()
			embedding, err := d.llmClient.Embed(ctx, []map[string]interface{}{
				{
					"type": "text",
					"text": data,
				},
			})
			if err != nil {
				log.Errorf("llmClient.Embed failed, err: %v", err)
				return
			}
			obj := map[string]interface{}{
				"content": data,
			}
			// save data to weaviate
			if err := dao.Create(ctx, d.class, obj, embedding); err != nil {
				log.Errorf("dao.Create failed, err: %v", err)
			}
		}(i, data)
	}
	wg.Wait()
	return nil
}

// Query approximation data from weaviate, generate a whole prompt send to llm generate answer
func (d *TextManager) Query(ctx context.Context, prompt string) (string, error) {
	log := logrus.WithContext(ctx)
	log.Infof("Query, Prompt: %s", prompt)
	log.Infof("generate embedding, prompt: %s", prompt)
	results := make([]string, 0)
	Embeddings, err := d.llmClient.Embed(ctx, []map[string]interface{}{
		{
			"type": "text",
			"text": prompt,
		},
	})
	if err != nil {
		log.Warnf("llmClient.Embed failed, err: %v", err)
	}
	if len(Embeddings) != 0 {
		log.Infof("query approximation data from weaviate")
		results, err = dao.Query(ctx, d.class, prompt, Embeddings)
		if err != nil {
			log.Warnf("dao.Query failed, err: %v", err)
		}
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
