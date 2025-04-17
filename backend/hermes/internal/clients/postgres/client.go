package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DataBaseURL string `json:"DATABASE_URL"`

	DataBaseMaxConns           int `json:"DATABASE_MAX_CONNS" cfgDefault:"4"`
	DataBaseMinConns           int `json:"DATABASE_MIN_CONNS" cfgDefault:"0"`
	DataBaseHealthCheckSeconds int `json:"DATABASE_HEALTHCHECK_SECONDS" cfgDefault:"60"`
}

func NewClient(cfg Config) (*Client, error) {
	config, err := pgxpool.ParseConfig(cfg.DataBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = int32(cfg.DataBaseMaxConns)
	config.MinConns = int32(cfg.DataBaseMinConns)
	config.HealthCheckPeriod = time.Duration(cfg.DataBaseHealthCheckSeconds) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Client{pool}, nil
}

type Client struct {
	*pgxpool.Pool
}

func (c *Client) Release() error {
	c.Close()
	return nil
}
