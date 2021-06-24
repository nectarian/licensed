package machineid

import (
	"bytes"
	"errors"
	"io/ioutil"
)

// -------------------------------------------------------------------------------------

const (
	machineIDInDBus = "/var/lib/dbus/machine-id"
	machineIDInEtc  = "/etc/machine-id"
)

// GetMachineID returns the machine-id specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// see https://www.man7.org/linux/man-pages/man5/machine-id.5.html
func GetMachineID() ([]byte, error) {
	if id, err := ioutil.ReadFile(machineIDInDBus); err == nil {
		return bytes.TrimSpace(id), err
	}
	if id, err := ioutil.ReadFile(machineIDInEtc); err == nil {
		return bytes.TrimSpace(id), err
	}
	return nil, errors.New("read machine-id failed")
}

// -------------------------------------------------------------------------------------
