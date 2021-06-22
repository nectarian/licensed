package plugin

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"time"
)

// Verify verify machine-id signature
func Verify(cc CreateClientIDFunc, expiration int64, sign []byte, publicKey *rsa.PublicKey) error {
	var err error
	// verify expiration if exists
	if expiration != 0 {
		n := time.Now().UnixNano()
		if n > expiration {
			return errors.New("license expiration")
		}
	}

	// get ClientID
	var clientID ClientID = cc()

	if clientID.Load() != nil {
		return fmt.Errorf("load clientID failed : %w", err)
	}

	var raw []byte
	if raw, err = Wrap(clientID, expiration); err != nil {
		return err
	}

	// verify sign
	hash := sha512.New()
	if _, err = hash.Write(raw); err != nil {
		log.Fatal("sha512sum(raw)) : ", err)
	}
	sum := hash.Sum(nil)
	if err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, sum, sign); err != nil {
		return fmt.Errorf("rsa verify failed %w", err)
	}
	return nil
}
