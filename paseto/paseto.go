package paseto

import (
	"time"

	"github.com/o1egl/paseto"

	"golang.org/x/crypto/ed25519"
)

type PasetoTokenGenerator struct {
	*paseto.V2
}

func NewPasetoTokenGenerator() (*PasetoTokenGenerator, ed25519.PublicKey, ed25519.PrivateKey, error) {

	publicKey, privateKey, err := ed25519.GenerateKey(nil)

	if err != nil {
		return nil, nil, nil, err
	}

	return &PasetoTokenGenerator{
		paseto.NewV2(),
	}, publicKey, privateKey, nil
}

func (p *PasetoTokenGenerator) GeneratePasetoToken(privateKey ed25519.PrivateKey, footer *string, expiration time.Time, claims map[string]string) (*string, error) {

	jsonToken := &paseto.JSONToken{
		Expiration: expiration,
	}

	// Add custom claim	to the token
	for k, v := range claims {
		jsonToken.Set(k, v)
	}

	// Sign data
	token, err := p.Sign(privateKey, jsonToken, footer)
	// token = "v2.public.eyJkYXRhIjoidGhpcyBpcyBhIHNpZ25lZCBtZXNzYWdlIiwiZXhwIjoiMjAxOC0wMy0xMlQxOTowODo1NCswMTowMCJ9Ojv0uXlUNXSFhR88KXb568LheLRdeGy2oILR3uyOM_-b7r7i_fX8aljFYUiF-MRr5IRHMBcWPtM0fmn9SOd6Aw.c29tZSBmb290ZXI"
	if err != nil {
		return nil, err
	}
	return &token, nil

}
