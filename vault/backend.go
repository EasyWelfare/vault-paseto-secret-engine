package vault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/EasyWelfare/vault-paseto-secret-engine/paseto"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type backend struct {
	*framework.Backend
	paseto  *paseto.PasetoTokenGenerator
	storage logical.Storage
	config  *Config
}

const (
	helpMessage      = "Paseto token generator"
	publicKeyConst   = "public"
	privateKeyConst  = "private"
	displayNameConst = "displayName"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	pasetoGenerator, publicKey, privateKey, err := paseto.NewPasetoTokenGenerator()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	b := &backend{
		paseto:  pasetoGenerator,
		storage: conf.StorageView,
	}

	b.storage.Put(ctx, &logical.StorageEntry{
		Key:      publicKeyConst,
		Value:    publicKey,
		SealWrap: false,
	})

	b.storage.Put(ctx, &logical.StorageEntry{
		Key:      privateKeyConst,
		Value:    privateKey,
		SealWrap: false,
	})

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(helpMessage),
		BackendType: logical.TypeLogical,
	}

	b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	b.Backend.Setup(ctx, conf)

	return b, nil
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "paseto/token",

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.generateToken,
					Summary:  "Generate Paseto token.",
				},
			},
		},
		{
			Pattern: "paseto/public",
			Fields:  map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.readPublicKey,
					Summary:  "Return Paseto public key.",
				},
			},
		},
		{
			Pattern: "paseto/config",
			Fields: map[string]*framework.FieldSchema{
				"footer": {
					Type:        framework.TypeString,
					Description: "Specify the paseto footer.",
				},
				"ttl": {
					Type:        framework.TypeInt,
					Description: "Specify the paseto token TTL in seconds.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.configurePasetoGenerator,
					Summary:  "Configure Paseto token generator.",
				},
			},
		},
	}
}

type Config struct {
	footer string
	ttl    time.Duration
}

func (b *backend) configurePasetoGenerator(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.config != nil {
		return logical.ErrorResponse("Config was already set."), nil
	}

	b.config = &Config{}
	b.config.footer = data.Get("footer").(string)
	b.config.ttl = time.Second * data.Get("ttl").(time.Duration)
	//TODO enrich response
	return nil, nil

}

func (b *backend) generateToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}
	claim := map[string]string{displayNameConst: req.DisplayName}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	return paseto.GeneratePasetoToken()

}
func (b *backend) readPublicKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}
	// retrieve and return publicKey
	return nil, nil
}
