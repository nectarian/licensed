package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// LoadPublicKey ready input path and parse it as RSA public key
func LoadPublicKey(f string) (*rsa.PublicKey, error) {
	var err error
	var block *pem.Block
	if block, err = ReadPEMBlock(f); err != nil {
		return nil, err
	}
	var key *rsa.PublicKey
	if key, err = x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
		log.Fatal("parse public key : ", err)
	}
	return key, err
}

// LoadPrivateKey read input path and parse it as RSA private key
func LoadPrivateKey(f string) (*rsa.PrivateKey, error) {
	var err error
	var block *pem.Block
	if block, err = ReadPEMBlock(f); err != nil {
		return nil, err
	}

	var key *rsa.PrivateKey
	if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		log.Fatal("parse private key : ", err)
	}
	return key, err
}

// WritePEMBlock build pem block and write to ouput file
func WritePEMBlock(outputPath string, ktype string, metadata map[string]string, raw []byte) error {
	var err error
	var folder = filepath.Dir(outputPath)
	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("create output folder error : %w", err)
	}

	var kFile *os.File
	if kFile, err = os.Create(outputPath); err != nil {
		return fmt.Errorf("create key file error : %w", err)
	}
	defer kFile.Close()

	block := pem.Block{Type: ktype, Bytes: raw}
	if metadata != nil {
		block.Headers = metadata
	}
	if err = pem.Encode(kFile, &block); err != nil {
		return fmt.Errorf("output key file error : %w", err)
	}
	return nil
}

// ReadPEMBlock read input path and decode as pem block
func ReadPEMBlock(inputPath string) (*pem.Block, error) {
	var err error
	var raw []byte
	if raw, err = ioutil.ReadFile(inputPath); err != nil {
		return nil, fmt.Errorf("read pem file failed : %v", err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("bad pem file : %v", err)
	}
	return block, nil
}
