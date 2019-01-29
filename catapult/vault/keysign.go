package vault

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/marcboudreau/go-devops-talk/catapult"
)

// KeySigningService is an implementation that uses Vault to handle the key signing.
type KeySigningService struct {
	catapult.KeySigningService

	role   string
	client *vaultapi.Client
}

// New creates a new KeySigningService instance that uses Vault to sign the provided
// keys.
func New(role string) (*KeySigningService, error) {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &KeySigningService{
		role:   role,
		client: client,
	}, nil
}

// SignKey signs the provided key using the underlying Vault client.
func (p *KeySigningService) SignKey(publicKey io.Reader, principal string) (io.Reader, error) {
	if publicKey == nil {
		return nil, errors.New("no publicKey reader provided to SignKey method")
	}

	if principal == "" {
		return nil, errors.New("no principal for certificate provided to SignKey method")
	}

	publicKeyBytes, err := ioutil.ReadAll(publicKey)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"public_key":       string(publicKeyBytes),
		"valid_principals": principal,
		"cert_type":        "user",
		"extensions": map[string]string{
			"permit-port-forwarding": "",
			"permit-pty":             "",
		},
	}

	secret, err := p.client.SSH().SignKey(p.role, data)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(secret.Data["signed_key"].(string)), nil
}
