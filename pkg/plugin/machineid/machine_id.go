package machineid

import (
	"github.com/nectarian/licensed/pkg/plugin"
)

// -------------------------------------------------------------------------------------

// MachineID machine-id
type MachineID struct {
	machineID string
}

// -------------------------------------------------------------------------------------

// CreateClientID create MachineID
func CreateClientID(opts ...plugin.Option) plugin.ClientID {
	m := MachineID{}
	for _, opt := range opts {
		opt(&m)
	}
	return &m
}

// Load load data
func (m *MachineID) Load() error {
	machineID, err := GetMachineID()
	if err != nil {
		return err
	}
	m.machineID = string(machineID)
	return nil
}

// Serialize serialize machine-id
func (m MachineID) Serialize() ([]byte, error) {
	return []byte(m.machineID), nil
}

// Deserialize deserialize machine-id
func (m *MachineID) Deserialize(raw []byte) error {
	m.machineID = string(raw)
	return nil
}
