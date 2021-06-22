package factory

import (
	"errors"

	"github.com/nectarian/licensed/pkg/plugin"
	"github.com/nectarian/licensed/pkg/plugin/baseboard"
	"github.com/nectarian/licensed/pkg/plugin/mac"
	"github.com/nectarian/licensed/pkg/plugin/machineid"
	"github.com/nectarian/licensed/pkg/plugin/plaintext"
)

// CreateClientIDFuncFactory CreateClientIDFunc factory
func CreateClientIDFuncFactory(name string) (plugin.CreateClientIDFunc, error) {
	switch name {
	case "machine-id":
		return machineid.CreateClientID, nil
	case "mac":
		return mac.CreateClientID, nil
	case "text":
		return plaintext.CreateClientID, nil
	case "baseboard":
		return baseboard.CreateClientID, nil
	default:
		return nil, errors.New("bad type")
	}
}
