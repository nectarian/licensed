package byted

import (
	"github.com/nectarian/licensed/pkg/plugin"
)

// -------------------------------------------------------------------------------------

// ByteClientID byte ClientID , just use for postmark do encrypt
type ByteClientID struct {
	RAW []byte
}

// -------------------------------------------------------------------------------------

// CreateClientID create byte client id
func CreateClientID(opts ...plugin.Option) plugin.ClientID {
	m := ByteClientID{}
	for _, opt := range opts {
		opt(&m)
	}
	return &m
}

// Load load data
func (m *ByteClientID) Load() error {
	return nil
}

// Serialize serialize byte client
func (m ByteClientID) Serialize() ([]byte, error) {
	return m.RAW, nil
}

// Deserialize deserialize byte client
func (m *ByteClientID) Deserialize(raw []byte) error {
	m.RAW = raw
	return nil
}
