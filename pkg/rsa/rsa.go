package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// RSAKeys manages the RSA private and public key pair for JWT operations.
type RSAKeys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// LoadOrCreate checks if key files exist on disk.
// If found, it loads them. Otherwise, it generates a new RSA keypair and saves it.
func LoadOrCreate(privPath, pubPath string, bitSize int) (*RSAKeys, error) {
	if fileExists(privPath) && fileExists(pubPath) {
		return LoadFromFiles(privPath, pubPath)
	}

	keys, err := Generate(bitSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new keys: %w", err)
	}

	if err := keys.SavePEMFiles(privPath, pubPath); err != nil {
		return nil, fmt.Errorf("failed to save generated keys: %w", err)
	}

	return keys, nil
}

// LoadFromFiles reads private and public PEM key files from disk.
func LoadFromFiles(privPath, pubPath string) (*RSAKeys, error) {
	privBytes, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &RSAKeys{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

// NewFromPEM parses a raw PEM-encoded private key string (e.g. loaded from .env).
func NewFromPEM(pemData string) (*RSAKeys, error) {
	cleanPEM := strings.ReplaceAll(pemData, "\\n", "\n")

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cleanPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key from PEM: %w", err)
	}

	return &RSAKeys{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// Generate creates a new in-memory RSA keypair of the specified bit size (e.g., 2048 or 4096).
func Generate(bitSize int) (*RSAKeys, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	return &RSAKeys{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// GetPrivateKey returns the internal RSA private key.
func (k *RSAKeys) GetPrivateKey() *rsa.PrivateKey {
	if k == nil {
		return nil
	}
	return k.privateKey
}

// GetPublicKey returns the internal RSA public key.
func (k *RSAKeys) GetPublicKey() *rsa.PublicKey {
	if k == nil {
		return nil
	}
	return k.publicKey
}

// SavePEMFiles exports the keypair to PEM-encoded files on disk.
func (k *RSAKeys) SavePEMFiles(privPath, pubPath string) error {
	if k == nil || k.privateKey == nil {
		return errors.New("cannot export uninitialized RSA keys")
	}

	// 1. Export Private Key (PKCS#1)
	privBytes := x509.MarshalPKCS1PrivateKey(k.privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	if err := os.WriteFile(privPath, pem.EncodeToMemory(privBlock), 0600); err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	// 2. Export Public Key (PKIX / PKCS#8)
	pubBytes, err := x509.MarshalPKIXPublicKey(k.publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	if err := os.WriteFile(pubPath, pem.EncodeToMemory(pubBlock), 0644); err != nil {
		return fmt.Errorf("failed to write public key file: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
