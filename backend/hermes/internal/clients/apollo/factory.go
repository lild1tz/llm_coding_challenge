package apollo

import (
	"context"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

type Client interface {
	PredictTableFromText(ctx context.Context, text string) (models.Table, error)
	PredictTableFromImage(ctx context.Context, image []byte) (models.Table, error)
	PredictTextFromAudio(ctx context.Context, audio []byte) (string, error)

	CheckVerbiage(ctx context.Context, text string) (bool, error)

	Release() error
}

type Config struct {
	ApolloURL string `json:"APOLLO_URL"`

	IsStub bool `json:"APOLLO_IS_STUB" cfgDefault:"true"`
}

func NewClient(cfg Config) Client {
	if cfg.IsStub {
		return &stubClient{}
	}

	return newClient(cfg)
}
