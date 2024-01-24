// Configuration
package cfg

import (
	"encoding/json"
	"fmt"
	"io"
)

type Map struct {
	// This is the original name of the account.
	OriginalName string `json:"original"`
	// Here is how to categorize it.
	Category string `json:"category"`
}

type Schema struct {
	ID         string `json:"account_id"`
	AccountMap []Map  `json:"account_map"`
}

// LoadSchema loads the JSON schema from the supplied reader.
func LoadSchema(r io.Reader) (*Schema, error) {
	d := json.NewDecoder(r)
	var s Schema
	if err := d.Decode(&s); err != nil {
		return nil, fmt.Errorf("while loading config schema: %w", err)
	}
	return &s, nil
}

type Instance struct {
	id string
	m  map[string]string
}

func (i Instance) GetID() string {
	return i.id
}

// GetCat returns a category for the given account. If not found, the account
// name is returned.
func (i Instance) GetCat(acc string) string {
	v, ok := i.m[acc]
	if !ok {
		return acc
	}
	return v
}

// New creates a new configuration based on the passed in config schema.
func New(s *Schema) *Instance {
	var i Instance
	i.id = s.ID
	i.m = map[string]string{}
	for _, m := range s.AccountMap {
		i.m[m.OriginalName] = m.Category
	}
	return &i
}
