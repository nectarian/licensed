package plugin

import (
	"encoding/binary"
	"errors"
)

// ClientID the client unique identification
type ClientID interface {
	// Load clientID data
	Load() error
	// Serialize parameter is unix timestamp of expiration
	Serialize() ([]byte, error)
	// Deserialize deserialize and return unix timestamp of expiration
	Deserialize([]byte) error
}

// CreateClientIDFunc The function signature of create ClientID
type CreateClientIDFunc func(...Option) ClientID

// Option for clientID
type Option func(ClientID)

func Wrap(client ClientID, expiration int64) ([]byte, error) {
	raw, err := client.Serialize()
	if err != nil {
		return nil, err
	}
	var msg = make([]byte, len(raw)+8+2)
	// write len
	binary.LittleEndian.PutUint16(msg[0:2], uint16(len(raw)))
	// write data
	copy(msg[2:len(raw)+2], raw)
	// write expiration
	binary.LittleEndian.PutUint64(msg[len(raw)+2:], uint64(expiration))
	return msg, nil
}

func Unwrap(raw []byte, ccf CreateClientIDFunc, opts ...Option) (ClientID, int64, error) {
	if len(raw) < 10 {
		return nil, 0, errors.New("bad input sequence , too short")
	}
	// read len
	var size = int(binary.LittleEndian.Uint16(raw[0:2]))
	// read data
	var data = raw[2 : 2+size]
	client := ccf(opts...)
	if err := client.Deserialize(data); err != nil {
		return nil, 0, err
	}
	// read expiration
	var expiration = int64(binary.LittleEndian.Uint64(raw[2+size:]))
	return client, expiration, nil
}
