package vault

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) readPublicKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Print("reading publicKey")
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}
	publicKeyStorageEntity, err := b.storage.Get(ctx, publicKeyConst)
	if err != nil {
		log.Printf("not found public key %s in storage: %v", privateKeyConst, err)
		return logical.ErrorResponse(err.Error()), err
	}
	response := &logical.Response{
		Data: map[string]interface{}{
			"publicKey": publicKeyStorageEntity.Value,
		},
	}

	return logical.RespondWithStatusCode(response, req, http.StatusOK)
}
