/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strings"
)

var keyFolder = "/etc/secrets/"
var privateKeyCache = make(map[string]*rsa.PrivateKey)
var publicKeyCache = make(map[string]*rsa.PublicKey)

func SetKeyFolder(path string) {
	path = strings.TrimSpace(path)
	if path != "" {
		c := path[len(path)-1]
		if c != '/' && c != '\\' {
			path += "/"
		}
	}
	keyFolder = path
}

func GetPublicKey(key string) (*rsa.PublicKey, error) {
	if res, ok := publicKeyCache[key]; ok {
		return res, nil
	}
	publicKey, err := LoadPublicKey(key)
	if err != nil {
		return nil, err
	}
	publicKeyCache[key] = publicKey
	return publicKey, nil
}

func GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	if res, ok := privateKeyCache[key]; ok {
		return res, nil
	}
	privateKey, err := LoadPrivateKey(key)
	if err != nil {
		return nil, err
	}
	privateKeyCache[key] = privateKey
	return privateKey, nil
}

func LoadEncodedData(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodeByteLine(data)
}

func SaveEncodedData(path string, data []byte) error {
	data = EncodeByteLine(data)
	return ioutil.WriteFile(path, data, 0644)
}

func loadByteKey(key string) ([]byte, error) {
	return LoadEncodedData(keyFolder + key)
}

func saveByteKey(key string, data []byte) error {
	return SaveEncodedData(keyFolder+key, data)
}

func LoadPublicKey(key string) (*rsa.PublicKey, error) {
	data, err := loadByteKey(key)
	if err != nil {
		return nil, err
	}
	res, err1 := x509.ParsePKCS1PublicKey(data)
	if err1 != nil {
		return nil, err1
	}
	return res, nil
}

func LoadPrivateKey(key string) (*rsa.PrivateKey, error) {
	data, err := loadByteKey(key)
	if err != nil {
		return nil, err
	}
	res, err1 := x509.ParsePKCS8PrivateKey(data)
	if err1 != nil {
		return nil, err1
	}
	privKey, ok := res.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Key " + key + " is not RSA private key")
	}
	return privKey, nil
}

func SavePublicKey(keyName string, key *rsa.PublicKey) error {
	data := x509.MarshalPKCS1PublicKey(key)
	return saveByteKey(keyName, data)
}

func SavePrivateKey(keyName string, key *rsa.PrivateKey) error {
	data := x509.MarshalPKCS1PrivateKey(key)
	return saveByteKey(keyName, data)
}

func CreatePublicPrivatePair(publicName string, privateName string, bits int) error {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	err = SavePrivateKey(privateName, privKey)
	if err != nil {
		return err
	}
	pubKey := privKey.Public()
	rsaPubKey, ok := pubKey.(rsa.PublicKey)
	if !ok {
		return errors.New("Public key is not RSA")
	}
	return SavePublicKey(publicName, &rsaPubKey)
}
