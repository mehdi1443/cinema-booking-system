package rsa

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type RSAKeys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func ParsePrivateKey(pemBytes []byte) (*RSAKeys, error) {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &RSAKeys{
		privateKey: parsedKey,
		publicKey:  &parsedKey.PublicKey,
	}, nil
}

func (k *RSAKeys) GetPrivateKey() *rsa.PrivateKey {
	if k == nil {
		return nil
	}
	return k.privateKey
}

func (k *RSAKeys) GetPublicKey() *rsa.PublicKey {
	if k == nil {
		return nil
	}
	return k.publicKey
}
