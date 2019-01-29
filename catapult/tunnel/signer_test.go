package tunnel

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These are dummy SSH keys that were generated specifically for testing this code.

var testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAzleDOIZvyJMScUb2nI3jWP5TkHWsccvSONq9I8B6kF5NTa9o
CwFD8oBpuel/Vceaj1ZllBD0z2aM3KLH11mRuytj+X4hCScS2bGwMW4gIfbhWHP/
UelRkYXjpIDaVa0sC7bvP5IqfTU3XzvfzQ0n5txCbJTz9P2HCLOSCr/vvWVE0C+0
5Uuh8ZLo1yR9XjCaoHeWz9mLkAd8N4SLVKKevNm9ANohKHqHseVWS8EYORS8cK5l
uwGtVklaCeM8Y2yeqlN6uos44BPed6GIFEUH1NT5nYo3VNACMmxuTvo5Ke9pnjbD
8u1ChdIJ9+J9v8ib5rlnxsZLKGnE/i7xKyO54QIDAQABAoIBAQC8t3cxJGtqM3DD
n5Z/Opn606hTz/vmm/Zpv00LPMgb79OdwFZU8lRVnzKTMUYfiw4GGIuQry1n3q/o
PwytHaNWXunxznSibVUlOwkdPE4xIh2Zi4WxQcYzZRP9aUfG4joNgMMyDhnEJ/67
oAQhAu5Ci6JIsraok5OD1tAz+rVmLml62W7gWYvRrmbAqlXXNvBR2xNhuoSC6kIy
LN5dlLYGvI86LrPkQRcPvuE3CIBDoTPlNu+6qnLltae+p9jcPa/UgNxeaLyxl3Oa
C+TvG6kdVwnxL0haarvfeq0PEGCE9YpZPa3KN6PfRkvGRbYPZE8171M1oKWSYuca
rArKBDm9AoGBAPtG2ozBiHblcQLtM8c8Omy7nujnEWaht8C+LykpI6Me9rdf6fyL
0yHLm2u8egrK8uoVMV3s9kWQu4T5xhFCSxLPzmh/HoANd3Q+3cMu2omXqGSzJi5o
f8WXxK0lLxlE01F123zIRAQl1x2zCcLbiQis1wMPrGtWdfZIN6zbhuoTAoGBANI4
bnsUT/T75tHUzXx9IJRefvCtxzvHaX6cOAbKLS6Sfld2WQ0ubYL/p43T8Zr6u1Lk
QEenZCMu0bjoOXYyQE0xW/ni35g2Tvwn3Lqevs5gOtzdMlf44wwSXXQqp4KrTSd1
E5oPD1zzkYS7pIP6kkJe346N4u5iVgjIAHF6Bgq7AoGBAKG5JIhbNz1uxqGfkSe1
99RrnQdBUM3BX8bJoQjY1XrzPs8fCDXmuGiT5uAcWl//5wAJy9Ar5wU29bnMGFKb
XD4rSSmwRy0bfbpvi8NHsJfF6DeHphdQYowF9iuKNxoIVgmj1TQmoMAaqq0OwkWL
jlLrCyeJOuuKpjlwmYTDdb3JAoGBAJxsvVj8TlrfLmwoyxa9DQcaIZ750Gyc/9Tk
bZQv0Nr8yuJOAAmc6IQ3s/gHI5rMw6L0kRhAaHT9m7TZqBhZYBuQhP42YWaj0rYy
+z4qbZSnamV6esGXQ2tyJvQP4UGMMgcQSRuz5RynaTq5XbuPMlIMwpES5y+3IIm8
OQg3YlONAoGAfWFmipVXPWYSQacWl8j8AB9Vv2E+DzYhZ+WVKGvbW2Goq4SL3qbq
42fx2nVMACaiGClFhDHqMMsl65EzXDjztSP8JJJtnpQxESi33a3dtlJhv8mbjog4
tDUPTmrV0gvK3iCU4ouRImvGer3ZxzItOf9guWYG880mL/5HvlC3Gdk=
-----END RSA PRIVATE KEY-----`

var testAlternatePrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAs8M8dDP7BnoqXFIzogW9GXSksf5qLmjBTxGaKiPh51SSXoun
kouYp/h43ATLNXp1gpjPu/6wN4JP6CICmYHeow6TFqgH6QDQTcX5eBEnLilHTjI6
Dm8OceYRkKZvnsiZ1nX0bPr2Hoy2gnBcGeZ0Iej0Zp0EZsy5hZGQR3XsAGAGv5fm
F2ivXCEfBjADj/aOnQN9qP+y49oFDQ+BzxyP7HFGLtmxxz8NxR856ynWtKE7Tg0a
ae/1OMBzvwhWohR9RQPSyJp9pRyBtZ/tQl8Y8gcCmOUEmmfm6SnGVqEkdLVSoUHJ
plRgPDwHR6Dcl0pQ2Oc7jAIvFPegeEWWinUvtQIDAQABAoIBACV7YwKTyBasqGKB
nVR+P9Sr9p9KfhdZLl1vPSbmPnc+shpWokUzKEx1ybWOplRrSU9Gz6HSVCnu2Px9
Au/BHYwAQPkrZiLWUZi12/OGGBZO8xhB7ssNqwRixWzU11dTWohWJfYHSgRKsM93
4CxiWfqsGTHAMafBrjlEhcHwu5nEEo5WStv2uT1nP/nQoF2NbWqc8tELtYQmHbPT
Dw2U9B8hnBRkQBF0KU86S6fNc6TR2VSALukasPNpdo+khmNyPBc/ewXiUFH/WggB
wfvwHdcqXtfrtbgYsNKj7ZT/cNqhxuaO/oYcNgg5Ww7DuWZMpbBxgHsvEYQoiNic
1nv2he0CgYEA1tjA4HzMEAVNCLWiHs0aA8Lc3YqeTfrdqWrAeaJAH2S32Ye0z1bY
7Kz83+x0vgFTS6BKu9EnZecufMIImdGnicqrlPet+9X0jQVJDHjnzxQyXqa+tqJG
3HyuaoUzjcErDypEi1jfnlWdw2Bi1Y6nFkHqL80JFPNxduFSyfdocn8CgYEA1jIY
4w29cMhnSrGVxviE/8Bigo55jO22sGNBK8uYFbTI23bl2tHi5u/J34KhUuIS8KpD
3rHIM+6qacT7k/l/YBHfjVD7rYKAsr0LguqpJNL6zSe8uf46igiUAiiYz47A4cLA
YAWUJyZhp8817mq1jQlRhNRk+l9h/K53K/PxG8sCgYApEQ91GYWr/kdmRcmyV4QE
egfbtPZjc3NRQd1+onvdRFQ8GUt/YlteigZgCwOZmglA6GfAlM8SFGl3YWNhe4ip
tvqrI9i2zYPMPNlkr+unUnX6T6ceo9AlrxNruwBKtUS0xmCJvjgoPLdGNDNQHM9l
Wr4X2vpEbfAfSByaDFeDowKBgDJ9LJCYtIbZNj/NDFkSS4ddedr8anplriJ1wu7n
1rmHG7FfnV6vCqUU8KjWyeAXmg9Qkx7zGRXktBaqcAK5VOT4UZGl/S0xDdUT6pq5
ZZVghe3F6B6PZo0S6JB5sUt7gsu0dFQt/HYd/fboSPYiARl1kNmWlxbNVPOSzFR3
8+VxAoGAEmniUwDoU3aXQQt4dySG4crEg48VfgxihrhkonPgU7uP4f7wDj1ReTjc
KGHwmTYoGaqrQg5fOR4mmOGKiEobkp2X1lqPksR0KCe4NIVYs8jwcPe/FNV5r8bZ
B185BhoruS7h/vVP9O1q7XRCUheUHyjKEP4niBJdK9o9zYHh4Oo=
-----END RSA PRIVATE KEY-----`

var testPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOV4M4hm/IkxJxRvacjeNY/lOQdaxxy9I42r0jwHqQXk1Nr2gLAUPygGm56X9Vx5qPVmWUEPTPZozcosfXWZG7K2P5fiEJJxLZsbAxbiAh9uFYc/9R6VGRheOkgNpVrSwLtu8/kip9NTdfO9/NDSfm3EJslPP0/YcIs5IKv++9ZUTQL7TlS6HxkujXJH1eMJqgd5bP2YuQB3w3hItUop682b0A2iEoeoex5VZLwRg5FLxwrmW7Aa1WSVoJ4zxjbJ6qU3q6izjgE953oYgURQfU1PmdijdU0AIybG5O+jkp72meNsPy7UKF0gn34n2/yJvmuWfGxksoacT+LvErI7nh root@1ce7436f2c86`

var testSignedPublicKey = `ssh-rsa-cert-v01@openssh.com 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`

func TestCreateSigner(t *testing.T) {
	testcases := []struct {
		privateKey  io.Reader
		certificate io.Reader

		fails bool
	}{
		// privateKey is nil
		{
			certificate: strings.NewReader(testSignedPublicKey),
			fails:       true,
		},
		// certificate is nil
		{
			privateKey: strings.NewReader(testPrivateKey),
			fails:      true,
		},
		// improperly formatted private key
		{
			privateKey:  strings.NewReader(strings.Replace(testPrivateKey, "\n", "", -1)),
			certificate: strings.NewReader(testSignedPublicKey),
			fails:       true,
		},
		// Invalid certificate
		{
			privateKey:  strings.NewReader(testPrivateKey),
			certificate: strings.NewReader(testSignedPublicKey[:len(testSignedPublicKey)-5]),
			fails:       true,
		},
		// Mismatched private key and signed public key
		{
			privateKey:  strings.NewReader(testAlternatePrivateKey),
			certificate: strings.NewReader(testSignedPublicKey),
			fails:       true,
		},
		// IO Error while reading private key
		{
			privateKey:  &failingReader{},
			certificate: strings.NewReader(testSignedPublicKey),
			fails:       true,
		},
		// IO Error while reading certificate
		{
			privateKey:  strings.NewReader(testPrivateKey),
			certificate: &failingReader{},
			fails:       true,
		},
		// Passes
		{
			privateKey:  strings.NewReader(testPrivateKey),
			certificate: strings.NewReader(testSignedPublicKey),
		},
	}

	for _, testcase := range testcases {
		if !testcase.fails {
			fmt.Println()
		}

		signer, err := CreateSigner(testcase.privateKey, testcase.certificate)
		if testcase.fails {
			assert.NotNil(t, err)
			assert.Nil(t, signer)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, signer)
		}
	}
}

type failingReader struct {
	io.Reader
}

func (p *failingReader) Read(buffer []byte) (int, error) {
	return 1, errors.New("Forced error for testing")
}
