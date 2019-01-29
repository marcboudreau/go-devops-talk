package tunnel

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// CreateSigner creates an ssh.Signer instance with the provided private key
// and certificate (signed public key).
func CreateSigner(privateKey, certificate io.Reader) (ssh.Signer, error) {
	if privateKey == nil {
		return nil, errors.New("no private key reader provided")
	}

	if certificate == nil {
		return nil, errors.New("no signed public key reader provided")
	}

	privateKeyBytes, err := ioutil.ReadAll(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error encountered reading private key from reader: %s", err)
	}

	privateKeySigner, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %s", err)
	}

	certificateBytes, err := ioutil.ReadAll(certificate)
	if err != nil {
		return nil, fmt.Errorf("error encountered reading signed public key from reader: %s", err)
	}

	certificatePublicKey, _, _, _, err := ssh.ParseAuthorizedKey(certificateBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing signed public key: %s", err)
	}

	return ssh.NewCertSigner(certificatePublicKey.(*ssh.Certificate), privateKeySigner)

	// sshConfig := &ssh.ClientConfig{
	// 	User:            username,
	// 	Auth:            []ssh.AuthMethod{},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }

}
