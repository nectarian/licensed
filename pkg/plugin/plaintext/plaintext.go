package plaintext

import (
	"github.com/nectarian/licensed/pkg/plugin"
)

// -------------------------------------------------------------------------------------

// Plaintext used for do not validate client
type Plaintext struct {
	text string
}

// TEXT text value "license"
const TEXT = "license"

// -------------------------------------------------------------------------------------

// CreateClientID create plaintext client id
func CreateClientID(opts ...plugin.Option) plugin.ClientID {
	m := Plaintext{}
	for _, opt := range opts {
		opt(&m)
	}
	return &m
}

// Load load data
func (m *Plaintext) Load() error {
	if m.text == "" {
		m.text = TEXT
	}
	return nil
}

// Serialize serialize plaintext client id
func (m Plaintext) Serialize() ([]byte, error) {
	return []byte(m.text), nil
}

// Deserialize deserialize plaintext client id
func (m *Plaintext) Deserialize(s []byte) error {
	m.text = string(s)
	return nil
}

// -------------------------------------------------------------------------------------

// Text return a option to change default text
func Text(text string) plugin.Option {
	return func(ci plugin.ClientID) {
		t, ok := ci.(*Plaintext)
		if ok {
			t.text = text
		}
	}
}
