package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	log.Printf("here")
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
			Pattern:         "info",
			HelpSynopsis:    helpMessage,
			HelpDescription: helpMessage,
			Fields:          map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{},
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

type Config struct {
	Footer string        `json:"footer"`
	Ttl    time.Duration `json:"ttl"`
}

func (b *backend) configurePasetoGenerator(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("received configurePasetoGenerator call with req %v and data %v ", *req, *data)
	if b.config != nil {
		return logical.ErrorResponse("Config was already set."), nil
	}

	b.config = &Config{}
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
	config, _ := json.Marshal(*b.config)
	response := &logical.Response{
		Data: map[string]interface{}{
			"config": string(config),
		},
	}

	return logical.RespondWithStatusCode(response, req, http.StatusOK)

}
func (b *backend) generateToken(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
			"token": token,
		},
	}

	return logical.RespondWithStatusCode(response, req, http.StatusOK)

}
func (b *backend) readPublicKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
