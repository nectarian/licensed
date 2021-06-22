package client

import (
	"context"
	"crypto/rsa"
	"encoding/pem"
	"log"
	"time"

	"github.com/nectarian/licensed/pkg/plugin"
	"github.com/nectarian/licensed/pkg/plugin/factory"
	rsautil "github.com/nectarian/licensed/pkg/utils/rsa"
)

// WatchLicense run a goroutine and validate license
// ltype can be one of : machine-id mac text
// ltype should equal with license type
func WatchLicense(ctx context.Context, ltype, licensePath, publicKeyPath string) {
	// read license file
	var err error

	var block *pem.Block
	if block, err = rsautil.ReadPEMBlock(licensePath); err != nil {
		log.Fatal(err)
	}

	// read public key
	var publicKey *rsa.PublicKey
	if publicKey, err = rsautil.LoadPublicKey(publicKeyPath); err != nil {
		log.Fatal(err)
	}

	// build cleartext
	var expiration int64
	if block.Headers["expiration"] != "" {
		var expirationTime time.Time
		if expirationTime, err = time.Parse("2006-01-02", block.Headers["expiration"]); err != nil {
			log.Fatal("bad expiration fomat", err)
		}
		expiration = expirationTime.UnixNano()
	}

	var CreateClientFunc plugin.CreateClientIDFunc
	CreateClientFunc, err = factory.CreateClientIDFuncFactory(ltype)
	if err != nil {
		log.Fatalf("bad type , type must be one of : mac machine-id text")
	}

	if err = plugin.Verify(CreateClientFunc, expiration, block.Bytes, publicKey); err != nil {
		log.Fatalf("verify license failed : %v", err)
	}
	// verify signature
	tic := time.NewTicker(time.Hour)
a:
	for {
		select {
		case <-ctx.Done():
			break a
		case <-tic.C:
			if err = plugin.Verify(CreateClientFunc, expiration, block.Bytes, publicKey); err != nil {
				log.Fatalf("verify license failed : %v", err)
			}
		}
	}
}
