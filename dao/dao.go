package dao

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"sync"
)

var client *weaviate.Client
var once sync.Once

type TextQueryResult struct {
	Data struct {
		Get struct {
			TextData []struct {
				Content string `json:"content"`
			} `json:"TextData"`
		} `json:"Get"`
	} `json:"data"`
}

func InitWeaviate(config weaviate.Config, class string) {
	once.Do(func() {
		var err error
		client, err = weaviate.NewClient(config)
		if err != nil {
			logrus.Fatalln(err)
		}
		if err := CreateSchemaIfNotExist(context.Background(), class); err != nil {
			logrus.Fatalln(err)
		}
	})
}

// CreateSchemaIfNotExist if schema not exist, create it
func CreateSchemaIfNotExist(ctx context.Context, class string) error {
	exist, err := client.Schema().ClassExistenceChecker().WithClassName(class).Do(ctx)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = client.Schema().ClassCreator().WithClass(&models.Class{
		Class:               class,
		Description:         "work log data",
		InvertedIndexConfig: nil,
		ModuleConfig:        nil,
		MultiTenancyConfig:  nil,
		Properties: []*models.Property{
			{
				DataType:          []string{"string"},
				Description:       "",
				IndexFilterable:   nil,
				IndexInverted:     nil,
				IndexRangeFilters: nil,
				IndexSearchable:   nil,
				ModuleConfig:      nil,
				Name:              "content",
				NestedProperties:  nil,
				Tokenization:      "",
			},
		},
		ReplicationConfig: nil,
		ShardingConfig:    nil,
		VectorConfig:      nil,
		VectorIndexConfig: nil,
		VectorIndexType:   "",
		Vectorizer:        "none",
	}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Create a data object and save to weaviate
func Create(ctx context.Context, class string, obj map[string]interface{}, embedding []float32) error {
	log := logrus.WithContext(ctx)
	log.Infof("Create, class: %s obj: %+v", class, obj)
	_, err := client.Data().Creator().
		WithClassName(class).
		WithProperties(obj).
		WithVector(embedding).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Query query data from weaviate
func Query(ctx context.Context, class string, query string, embedding []float32) ([]string, error) {
	log := logrus.WithContext(ctx)
	log.Infof("Query, class: %s query: %s", class, query)
	vectorBuilder := &graphql.NearVectorArgumentBuilder{}
	vectorBuilder.WithVector(embedding)
	result, err := client.GraphQL().Get().
		WithClassName(class).
		WithFields(graphql.Field{Name: "content"}).
		WithNearVector(vectorBuilder).
		WithLimit(5).
		Do(ctx)
	if err != nil {
		log.Errorf("GraphQL Get failed, err: %v", err)
		return nil, err
	}
	dataBytes, err := result.MarshalBinary()
	if err != nil {
		log.Errorf("result.MarshalBinary failed, err: %v", err)
		return nil, err
	}
	var queryResult TextQueryResult
	err = json.Unmarshal(dataBytes, &queryResult)
	if err != nil {
		log.Errorf("json.Unmarshal TextQueryResult failed, err: %v", err)
		return nil, err
	}
	var contents []string
	for _, item := range queryResult.Data.Get.TextData {
		contents = append(contents, item.Content)
	}
	return contents, nil
}
