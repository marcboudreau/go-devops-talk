package vault

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

var testPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOV4M4hm/IkxJxRvacjeNY/lOQdaxxy9I42r0jwHqQXk1Nr2gLAUPygGm56X9Vx5qPVmWUEPTPZozcosfXWZG7K2P5fiEJJxLZsbAxbiAh9uFYc/9R6VGRheOkgNpVrSwLtu8/kip9NTdfO9/NDSfm3EJslPP0/YcIs5IKv++9ZUTQL7TlS6HxkujXJH1eMJqgd5bP2YuQB3w3hItUop682b0A2iEoeoex5VZLwRg5FLxwrmW7Aa1WSVoJ4zxjbJ6qU3q6izjgE953oYgURQfU1PmdijdU0AIybG5O+jkp72meNsPy7UKF0gn34n2/yJvmuWfGxksoacT+LvErI7nh root@1ce7436f2c86`

var testSignedPublicKey = `ssh-rsa-cert-v01@openssh.com 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\n`

func TestSignKey(t *testing.T) {
	testcases := []struct {
		publicKey io.Reader
		principal string
		handler   http.HandlerFunc

		fails bool
	}{
		// No publicKey reader
		{
			publicKey: nil,
			fails:     true,
		},
		// No principal string
		{
			publicKey: strings.NewReader(testPublicKey),
			fails:     true,
		},
		// Problem reading publicKey reader
		{
			publicKey: &failingReader{},
			principal: "test",
			fails:     true,
		},
		// Permission denied from Vault
		{
			publicKey: strings.NewReader(testPublicKey),
			principal: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(403)
				w.Write([]byte("* permission denied"))
			},
			fails: true,
		},
		// Key signing succeeded
		{
			publicKey: strings.NewReader(testPublicKey),
			principal: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(fmt.Sprintf(`{
	"request_id": "ca1b0bcc-0e08-99f0-bf1d-22551662d027",
	"lease_id": "",
	"lease_duration": 0,
	"renewable": false,
	"data": {
		"serial_number": "dda93051ae06a644",
		"signed_key": "%s"
	},
	"warnings": null
}`, testSignedPublicKey)))
			},
		},
	}

	config := vaultapi.DefaultConfig()

	for _, testcase := range testcases {
		server := httptest.NewServer(testcase.handler)

		config.Address = server.URL

		client, err := vaultapi.NewClient(config)
		service := &KeySigningService{
			role:   "test",
			client: client,
		}

		certificate, err := service.SignKey(testcase.publicKey, testcase.principal)
		if testcase.fails {
			assert.NotNil(t, err)
			assert.Nil(t, certificate)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, certificate)
		}

		server.Close()
	}
}

type failingReader struct {
	io.Reader
}

func (p *failingReader) Read(buffer []byte) (int, error) {
	return 1, errors.New("Forced error for testing")
}
