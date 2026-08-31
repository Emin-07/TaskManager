package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTestKeys generates a throwaway RSA keypair and writes it to t.TempDir(),
// returning the paths to the private and public PEM files.
func writeTestKeys(t *testing.T) (privPath, pubPath string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	privPath = filepath.Join(t.TempDir(), "private.pem")
	require.NoError(t, os.WriteFile(privPath, privPEM, 0600))

	pubDER := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubDER})
	pubPath = filepath.Join(t.TempDir(), "public.pem")
	require.NoError(t, os.WriteFile(pubPath, pubPEM, 0600))

	return privPath, pubPath
}

func TestTokenServ_CreateTokenAndParse(t *testing.T) {
	priv, pub := writeTestKeys(t)
	ts := NewTokenService(priv, pub)

	token, err := ts.CreateToken("42", "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	data, err := ts.ParseFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, "42", data["id"])
	assert.Equal(t, "admin", data["role"])
}

func TestTokenServ_TokensDifferPerSubject(t *testing.T) {
	priv, pub := writeTestKeys(t)
	ts := NewTokenService(priv, pub)

	alice, err := ts.CreateToken("1", "user")
	require.NoError(t, err)
	bob, err := ts.CreateToken("2", "admin")
	require.NoError(t, err)

	assert.NotEqual(t, alice, bob)
}

func TestTokenServ_ParseRejectsInvalidToken(t *testing.T) {
	priv, pub := writeTestKeys(t)
	ts := NewTokenService(priv, pub)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")

	_, err := ts.ParseFromRequest(req)
	assert.Error(t, err)
}
