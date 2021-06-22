package mac

import (
	"testing"
	"time"

	"github.com/nectarian/licensed/pkg/plugin"
)

func TestMAC(t *testing.T) {
	var err error
	var oclient plugin.ClientID = CreateClientID()

	if oclient.Load() != nil {
		t.Errorf("load MachineID failed : %v", err)
		return
	}

	var wraped []byte
	var expiration = time.Now().UnixNano()

	if wraped, err = plugin.Wrap(oclient, expiration); err != nil {
		t.Errorf("wrap MachineID client : %v", err)
		return
	}

	var nclient plugin.ClientID
	var nExpiration int64
	if nclient, nExpiration, err = plugin.Unwrap(wraped, CreateClientID); err != nil {
		t.Errorf("unwrap MachineID client : %v", err)
		return
	}

	var otext, ntext *MAC
	var ok bool
	if otext, ok = oclient.(*MAC); !ok {
		t.Errorf("ptext is not a MachineID : %v", err)
		return
	}

	if ntext, ok = nclient.(*MAC); !ok {
		t.Errorf("nptext is not a MachineID : %v", err)
		return
	}

	if otext.MAC.String() != ntext.MAC.String() {
		t.Errorf("unequal MachineID : original : %#v new : %#v", otext.MAC.String(), ntext.MAC.String())
		return
	}
	t.Logf("MachineID : original : %#v new : %#v", otext.MAC.String(), ntext.MAC.String())

	if nExpiration != expiration {
		t.Errorf("unequal expiration value : original %v new : %v", expiration, nExpiration)
		return
	}
	t.Logf("expiration value : original %v new : %v", expiration, nExpiration)
}
