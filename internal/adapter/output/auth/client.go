package auth

import (
	"context"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
)

var _ port.IdentityProvider = (*Client)(nil)

type Config struct {
	Domain       string
	ClientID     string
	ClientSecret string
}

type Client struct {
	config Config
}

func NewClient(cfg Config) *Client {
	return &Client{config: cfg}
}

func (c *Client) Save(ctx context.Context, identity *entity.Identity) error {
	return nil
}
