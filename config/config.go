package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once
var Cfg *Config

type Server struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type Ollama struct {
	Url   string `json:"url"`
	Model string `json:"model"`
}

type Weaviate struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	Class  string `json:"class"`
}

type Config struct {
	Server   Server   `json:"server"`
	Ollama   Ollama   `json:"ollama"`
	Weaviate Weaviate `json:"weaviate"`
}

func Init(filename string) {
	once.Do(func() {
		viper.SetConfigFile(filename)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		Cfg = &Config{
			Server: Server{
				Name: viper.GetString("server.name"),
				Addr: viper.GetString("server.addr"),
			},
			Ollama: Ollama{
				Url:   viper.GetString("ollama.url"),
				Model: viper.GetString("ollama.model"),
			},
			Weaviate: Weaviate{
				Host:   viper.GetString("weaviate.host"),
				Schema: viper.GetString("weaviate.schema"),
				Class:  viper.GetString("weaviate.class"),
			},
		}
		logrus.Infof("init Cfg: %+v", Cfg)
	})
}
