package llm

import "context"

type ModelManager interface {
	// Embed generates embeddings from a model.
	Embed(ctx context.Context, datas []map[string]interface{}) ([]float32, error)
	// Generate generates a response for a given prompt.
	Generate(ctx context.Context, prompt string) (string, error)
}
