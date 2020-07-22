package bunq

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

type CryptoService struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

const privateKeyFile = "id_rsa"
const publicKeyFile = "id_rsa.pub"

func (service *CryptoService) GetPrivatePem() []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(service.PrivateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	return pem.EncodeToMemory(&privBlock)
}

func (service *CryptoService) GetPublicPem() ([]byte, error) {
	publicKey := service.PrivateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	publicKeyPem := pem.EncodeToMemory(&publicKeyBlock)

	return publicKeyPem, nil
}

func (service *CryptoService) generateNewKeyPair() error {
	privateKey, err := service.generateNewPrivateKey()
	if err != nil {
		return err
	}

	err = privateKey.Validate()
	if err != nil {
		return err
	}

	service.PrivateKey = privateKey
	service.PublicKey = privateKey.PublicKey

	//write private key file
	err = ioutil.WriteFile(privateKeyFile, service.GetPrivatePem(), 0600)
	if err != nil {
		return err
	}

	//write public key file
	publicKey, err := service.GetPublicPem()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(publicKeyFile, publicKey, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (service *CryptoService) generateNewPrivateKey() (*rsa.PrivateKey, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)

	if err != nil {
		return nil, err
	}

	return key, nil
}

//func (service *CryptoService) generatePublicKey() (string, error) {
//	publicKey := service.PrivateKey.PublicKey
//	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
//	if err != nil {
//		return nil, err
//	}
//
//	publicKeyBlock := pem.Block{
//		Type:    "PUBLIC KEY",
//		Headers: nil,
//		Bytes:   publicKeyDer,
//	}
//	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))
//
//
//	publicRsaKey, err := ssh.NewPublicKey(service.PrivateKey.Public())
//	if err != nil {
//		return nil, err
//	}
//
//	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
//
//	return pubKeyBytes, nil
//}

func (service *CryptoService) Init() error {
	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		err := service.generateNewKeyPair()
		if err != nil {
			return err
		}
	}

	k, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(k)
	if block == nil {
		return errors.New("empty pem block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	service.PrivateKey = privateKey
	service.PublicKey = privateKey.PublicKey

	return nil
}

func NewCryptoService() (*CryptoService, error) {
	crypto := CryptoService{}
	err := crypto.Init()
	if err != nil {
		return nil, err
	}

	return &crypto, nil
}
