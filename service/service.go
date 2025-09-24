package service

import (
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/pkg/llm"
	"sync"
)

var once sync.Once
var Srv *Services

type Services struct {
	TextManger DataManager
}

// Init services
func Init() {
	once.Do(func() {
		Srv = &Services{}
		llmClient := llm.NewOllamaManager(config.Cfg.Ollama.Model, config.Cfg.Ollama.Url)
		Srv.TextManger = NewTextManager(llmClient)
	})

}
