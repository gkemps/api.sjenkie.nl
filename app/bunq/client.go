package bunq

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const tokenFile = "token"

type Service struct {
	BaseUrl      string
	ApiKey       string
	Token        string
	SessionToken string
	UserId       float64

	cryptoService *CryptoService
}

func (service *Service) Init() error {
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		err := service.registerDevice()
		if err != nil {
			return err
		}
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return err
	}

	service.Token = string(token)

	if service.SessionToken == "" {
		err := service.createSessionToken()
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) DoRequest(request *http.Request) (*http.Response, error) {
	err := service.addSignatureHeader(request)
	if err != nil {
		return nil, err
	}

	service.addAuthenticationHeader(request)

	return http.DefaultClient.Do(request)
}

func (service *Service) getPrivateKey() string {
	return string(service.cryptoService.GetPrivatePem())
}

func (service *Service) getPublicKey() string {
	publicKey, err := service.cryptoService.GetPublicPem()
	if err != nil {
		return ""
	}

	return string(publicKey)
}

func (service *Service) registerDevice() error {
	token, err := service.Installation()
	if err != nil {
		return err
	}

	service.Token = token

	//write public key file
	err = ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()
	ipAddresses, _ := service.getIpAddresses()

	err = service.NewDevice(fmt.Sprintf("api.sjenkie.nl registered from %s", hostname), ipAddresses)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) getIpAddresses() ([]string, error) {
	resp, err := http.DefaultClient.Get("https://api.ipify.org/")
	if err != nil {
		return make([]string, 0), err
	}

	ipAddress, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]string, 0), err
	}

	return []string{string(ipAddress)}, nil
}

func (service *Service) createSessionToken() error {
	sessionToken, userId, err := service.CreateSession()
	if err != nil {
		return err
	}

	service.SessionToken = sessionToken
	service.UserId = userId

	return nil
}

func (service *Service) addAuthenticationHeader(r *http.Request) {
	if service.SessionToken != "" {
		r.Header.Set("X-Bunq-Client-Authentication", service.SessionToken)
	} else {
		r.Header.Set("X-Bunq-Client-Authentication", service.Token)
	}
}

func (service *Service) addSignatureHeader(r *http.Request) error {
	if r.Method != http.MethodPost {
		return nil
	}

	bodyReader, err := r.GetBody()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return err
	}

	h := sha256.New()
	_, err = h.Write(body)
	if err != nil {
		return err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, service.cryptoService.PrivateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return err
	}

	r.Header.Set("X-Bunq-Client-Signature", base64.StdEncoding.EncodeToString(signature))

	return nil
}

func NewService(baseUrl string, apiKey string) (*Service, error) {
	cryptoService, err := NewCryptoService()
	if err != nil {
		return nil, err
	}

	service := Service{
		ApiKey:        apiKey,
		BaseUrl:       baseUrl,
		cryptoService: cryptoService,
	}

	err = service.Init()
	if err != nil {
		return nil, err
	}

	return &service, nil
}
