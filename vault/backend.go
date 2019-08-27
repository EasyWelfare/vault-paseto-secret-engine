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
	paseto  paseto.PasetoTokenGenerator
	storage logical.Storage
	config  Config
}

const (
	helpMessage      = "Paseto token generator"
	publicKeyConst   = "public"
	privateKeyConst  = "private"
	displayNameConst = "displayName"
)

type Config struct {
	Footer string        `json:"footer"`
	Ttl    time.Duration `json:"ttl"`
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	log.Printf("Factory initializing")
	pasetoGenerator, publicKey, privateKey, err := paseto.NewPasetoTokenGenerator()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	b := &backend{
		paseto:  *pasetoGenerator,
		storage: conf.StorageView,
	}

	if entry, _ := b.storage.Get(ctx, publicKeyConst); entry == nil {
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
	}

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
	log.Print("setting paths")
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
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.configurePasetoGenerator,
					Summary:  "Configure Paseto token generator.",
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.getConfiguration,
					Summary:  "Get configuration.",
				},
			},
		},
	}
}
