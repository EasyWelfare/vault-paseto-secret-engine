package vault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) generateToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Print("creating token")
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}
	claim := map[string]string{displayNameConst: req.DisplayName}

	privateKeyStorageEntity, err := b.storage.Get(ctx, privateKeyConst)
	if err != nil {
		log.Printf("not found private key %s in storage: %v", privateKeyConst, err)
		return logical.ErrorResponse(err.Error()), err
	}

	expirationTime := time.Now().Add(b.config.Ttl)

	token, err := b.paseto.GeneratePasetoToken(privateKeyStorageEntity.Value, &b.config.Footer, expirationTime, claim)
	if err != nil {
		log.Printf("error generating paseto token with footer %s, expirationTime %v and claim %v: %v", b.config.Footer, expirationTime, claim, err)
		return logical.ErrorResponse(err.Error()), err
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			"token":      token,
			"expiration": expirationTime,
		},
	}

	return logical.RespondWithStatusCode(response, req, http.StatusOK)

}
