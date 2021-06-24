package baseboard

import (
	"errors"
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/nectarian/licensed/pkg/plugin"
)

type Baseboard struct {
	baseboard string
}

// CreateClientID create baseboard client id
func CreateClientID(opts ...plugin.Option) plugin.ClientID {
	m := Baseboard{}
	for _, opt := range opts {
		opt(&m)
	}
	return &m
}

func (b *Baseboard) Load() error {
	baseboard, err := ghw.Baseboard()
	if err != nil {
		return fmt.Errorf("load baseboard failed : %w", err)
	}
	if baseboard.SerialNumber == "" {
		return errors.New("baseboard serialNumber is empty")
	}
	b.baseboard = baseboard.SerialNumber
	return nil
}

// Serialize baseboard client id
func (b Baseboard) Serialize() ([]byte, error) {
	return []byte(b.baseboard), nil
}

// Deserialize baseboard client id
func (b *Baseboard) Deserialize(raw []byte) error {
	b.baseboard = string(raw)
	return nil
}
