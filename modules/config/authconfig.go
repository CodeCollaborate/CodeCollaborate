package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/CodeCollaborate/Server/utils"
)

func rsaConfigSetup(rsaPrivateKeyLocation, rsaPrivateKeyPassword string) (*rsa.PrivateKey, error) {
	if rsaPrivateKeyLocation == "" {
		utils.LogWarn("No RSA Key given, generating temp one", nil)
		return GenRSA(4096)
	}

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		utils.LogWarn("No RSA private key found, generating temp one", nil)
		return GenRSA(4096)
	}

	privPem, _ := pem.Decode(priv)
	var pemBytes []byte

	if privPem.Type != "RSA PRIVATE KEY" {
		utils.LogWarn("RSA private key is of the wrong type", utils.LogFields{
			"Pem Type": privPem.Type,
		})
	}

	if rsaPrivateKeyPassword != "" {
		pemBytes, err = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
	} else {
		pemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(pemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(pemBytes); err != nil { // note this returns type `interface{}`
			utils.LogError("Unable to parse RSA private key, generating a temp one", err, utils.LogFields{})
			return GenRSA(4096)
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	if privateKey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		utils.LogError("Unable to parse RSA key, generating a temp one", err, utils.LogFields{})
		return GenRSA(4096)
	}

	utils.LogInfo("Loaded RSA key from file", utils.LogFields{})
	return privateKey, nil
}

// GenRSA returns a new RSA key of bits length
func GenRSA(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	utils.LogFatal("Failed to generate signing key", err, nil)
	return key, err
}
