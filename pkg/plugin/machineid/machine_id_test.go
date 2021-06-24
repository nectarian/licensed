package machineid

import (
	"testing"
	"time"

	"github.com/nectarian/licensed/pkg/plugin"
)

func TestMachineID(t *testing.T) {
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

	var otext, ntext *MachineID
	var ok bool
	if otext, ok = oclient.(*MachineID); !ok {
		t.Errorf("ptext is not a MachineID : %v", err)
		return
	}

	if ntext, ok = nclient.(*MachineID); !ok {
		t.Errorf("nptext is not a MachineID : %v", err)
		return
	}

	if otext.machineID != ntext.machineID {
		t.Errorf("unequal MachineID : original : %#v new : %#v", otext.machineID, ntext.machineID)
		return
	}
	t.Logf("MachineID : original : %#v new : %#v", *&otext.machineID, ntext.machineID)

	if nExpiration != expiration {
		t.Errorf("unequal expiration value : original %v new : %v", expiration, nExpiration)
		return
	}
	t.Logf("expiration value : original %v new : %v", expiration, nExpiration)
}
