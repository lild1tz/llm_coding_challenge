package apollo

import (
	"context"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

type Client interface {
	PredictTableFromText(ctx context.Context, text string) (models.Table, error)
}

type Config struct {
	ApolloURL string `json:"APOLLO_URL"`

	IsStub bool `json:"IS_STUB" cfgDefault:"true"`
}

func NewClient(cfg Config) Client {
	if cfg.IsStub {
		return &stubClient{}
	}

	return newClient(cfg)
}
