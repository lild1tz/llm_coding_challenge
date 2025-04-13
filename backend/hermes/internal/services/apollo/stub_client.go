package apollo

import (
	"context"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

type stubClient struct {
}

func (c *stubClient) PredictTableFromText(ctx context.Context, text string) (models.Table, error) {
	return models.Table{
		{
			Date:         time.Now().Format("2006-01-02"),
			Division:     "АОР",
			Operation:    "Внесение минеральных удобрений",
			Culture:      "Пшеница озимая товарная",
			PerDay:       117,
			PerOperation: 7381,
			ValDay:       1560,
			ValBeginning: 0,
		},
	}, nil
}

var x int

func (c *stubClient) CheckVerbiage(ctx context.Context, text string) (bool, error) {
	x++
	if x%2 == 0 {
		return true, nil
	}

	return false, nil
}
