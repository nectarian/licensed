package plaintext

import (
	"testing"
	"time"

	"github.com/nectarian/licensed/pkg/plugin"
)

func TestPlaintext(t *testing.T) {
	var err error
	var oclient plugin.ClientID = CreateClientID()

	if oclient.Load() != nil {
		t.Errorf("load Plaintext failed : %v", err)
		return
	}

	var wraped []byte
	var expiration = time.Now().UnixNano()

	if wraped, err = plugin.Wrap(oclient, expiration); err != nil {
		t.Errorf("wrap Plaintext client : %v", err)
		return
	}

	var nclient plugin.ClientID
	var nExpiration int64
	if nclient, nExpiration, err = plugin.Unwrap(wraped, CreateClientID); err != nil {
		t.Errorf("unwrap Plaintext client : %v", err)
		return
	}

	var otext, ntext *Plaintext
	var ok bool
	if otext, ok = oclient.(*Plaintext); !ok {
		t.Errorf("ptext is not a Plaintext : %v", err)
		return
	}

	if ntext, ok = nclient.(*Plaintext); !ok {
		t.Errorf("nptext is not a Plaintext : %v", err)
		return
	}

	if otext.text != ntext.text {
		t.Errorf("unequal Plaintext : original : %#v new : %#v", otext.text, ntext.text)
		return
	}
	t.Logf("Plaintext : original : %#v new : %#v", *&otext.text, ntext.text)

	if nExpiration != expiration {
		t.Errorf("unequal expiration value : original %v new : %v", expiration, nExpiration)
		return
	}
	t.Logf("expiration value : original %v new : %v", expiration, nExpiration)
}
