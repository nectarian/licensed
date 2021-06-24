package machineid

import (
	"golang.org/x/sys/windows/registry"
)

// GetMachineID get machine-id from windows registry
func GetMachineID() ([]byte, error) {
	var err error
	var key registry.Key
	if key, err = registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY); err != nil {
		return nil, err
	}
	s, _, err := key.GetStringValue("MachineGuid")
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}
