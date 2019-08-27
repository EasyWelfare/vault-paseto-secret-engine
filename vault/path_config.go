package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) configurePasetoGenerator(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("received configurePasetoGenerator call with req %v and data %v ", *req, *data)

	b.config = Config{}
	b.config.Footer = data.Get("footer").(string)
	ttl, ok := data.Get("ttl").(int)
	if !ok {
		return logical.ErrorResponse("tll is not an int"), nil
	}
	b.config.Ttl = time.Second * time.Duration(ttl)
	//TODO enrich response
	return nil, nil

}
func (b *backend) getConfiguration(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}
	log.Printf("footer: %s, ttl: %v", b.config.Footer, b.config.Ttl)
	config, _ := json.Marshal(b.config)
	response := &logical.Response{
		Data: map[string]interface{}{
			"config": string(config),
		},
	}

	return logical.RespondWithStatusCode(response, req, http.StatusOK)

}
