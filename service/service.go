package service

import (
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/pkg/llm/impl"
	"os"
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
		apiKey := config.Cfg.DouBao.ApiKey
		if len(apiKey) < 5 {
			apiKey = os.Getenv("API_KEY")
		}
		llmClient := impl.NewDouBao(apiKey)
		Srv.TextManger = NewTextManager(llmClient)
	})

}
