package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var once sync.Once
var Cfg *Config

type Server struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type DouBao struct {
	ApiKey string
}

type Weaviate struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	Class  string `json:"class"`
}

type Config struct {
	Server   Server   `json:"server"`
	DouBao   DouBao   `json:"doubao"`
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
			DouBao: DouBao{
				ApiKey: firstNonEmpty(os.Getenv("DOUBAO_API_KEY"), viper.GetString("doubao.apiKey")),
			},
			Weaviate: Weaviate{
				Host:   firstNonEmpty(os.Getenv("WEAVIATE_HOST"), viper.GetString("weaviate.host")),
				Schema: viper.GetString("weaviate.schema"),
				Class:  viper.GetString("weaviate.class"),
			},
		}
		logrus.Infof("init Cfg: %+v", Cfg)
	})
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
