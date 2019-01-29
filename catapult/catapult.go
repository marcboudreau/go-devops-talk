package catapult

import "io"

// KeySigningService is an interface defining methods for signing SSH keys.
type KeySigningService interface {
	SignKey(publicKey io.Reader, principal string) (io.Reader, error)
}
