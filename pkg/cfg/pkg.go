// Configuration
package cfg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Map struct {
	// This is the original name of the account.
	OriginalName string `json:"original"`
	// Here is how to categorize it.
	Category string `json:"category"`
	// ID is the unique account ID. If not specified, one will be generated.
	ID string `json:"id"`
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

type CatId struct {
	Category, Id string
}

type Instance struct {
	id string
	m  map[string]CatId
}

func (i Instance) GetID() string {
	return i.id
}

// GetCat returns a category for the given account. If not found, the account
// name is returned.
func (i Instance) GetCat(acc string) string {
	v := i.m[acc].Category
	if v == "" {
		// Log a warning since we could potentially add this into the config file.
		log.Printf("account has no category: %v", acc)
		return acc
	}
	return v
}

func (i Instance) GetAccID(acc string) string {
	v := i.m[acc].Id
	if v == "" {
		// Log a warning since we could potentially add this into the config file.
		log.Printf("account has no ID: %v", acc)
	}
	return v
}

// New creates a new configuration based on the passed in config schema.
func New(s *Schema) *Instance {
	var i Instance
	i.id = s.ID
	i.m = map[string]CatId{}
	for _, m := range s.AccountMap {
		i.m[m.OriginalName] = CatId{m.Category, m.ID}
	}
	return &i
}
