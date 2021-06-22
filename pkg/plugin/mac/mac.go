package mac

import (
	"fmt"
	"net"

	"github.com/nectarian/licensed/pkg/plugin"
)

// MAC Address
type MAC struct {
	name string
	MAC  net.HardwareAddr
	size uint16
}

// CreateClientID create MAC
func CreateClientID(opts ...plugin.Option) plugin.ClientID {
	m := MAC{}
	for _, opt := range opts {
		opt(&m)
	}
	return &m
}

// Load load data
func (m *MAC) Load() error {
	var err error
	var ifaces []net.Interface
	if ifaces, err = net.Interfaces(); err != nil {
		return fmt.Errorf("get network interface failed : %w ", err)
	}
	for _, v := range ifaces {
		if v.Flags&net.FlagUp == 0 {
			continue
		}
		if v.Flags&net.FlagLoopback != 0 {
			continue
		}
		if m.name == "" {
			m.MAC = v.HardwareAddr
			m.size = uint16(len(m.MAC))
			return nil
		}
		if m.name == v.Name {
			m.MAC = v.HardwareAddr
			m.size = uint16(len(m.MAC))
			return nil
		}

	}
	return fmt.Errorf("no network interface")
}

// Serialize mac client id
func (m MAC) Serialize() ([]byte, error) {
	return []byte(m.MAC), nil
}

// Deserialize mac client id
func (m *MAC) Deserialize(bslice []byte) error {
	m.MAC = bslice
	return nil
}

// NetworkName return a option to Specify network card name
func NetworkName(name string) plugin.Option {
	return func(ci plugin.ClientID) {
		mac, ok := ci.(*MAC)
		if ok {
			mac.name = name
		}
	}
}
