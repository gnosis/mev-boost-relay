package server

import (
	"sync"

	"github.com/flashbots/go-boost-utils/types"
)

type Datastore interface {
	GetValidatorRegistration(proposerPubkey types.PublicKey) (*types.SignedValidatorRegistration, error)
	SaveValidatorRegistration(entry types.SignedValidatorRegistration) error
	SaveValidatorRegistrations(entries []types.SignedValidatorRegistration) error
}

type MemoryDatastore struct {
	entries map[types.PublicKey]*types.SignedValidatorRegistration
	mu      sync.RWMutex

	// Used to count each request made to the datastore for each method
	requestCount map[string]int
}

// GetValidatorRegistration returns the validator registration for the given proposerPubkey. If not found then it returns (nil, nil). If
// there's a datastore error, then an error will be returned.
func (ds *MemoryDatastore) GetValidatorRegistration(proposerPubkey types.PublicKey) (*types.SignedValidatorRegistration, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	ds.requestCount["GetValidatorRegistration"]++
	return ds.entries[proposerPubkey], nil
}

func (ds *MemoryDatastore) SaveValidatorRegistration(entry types.SignedValidatorRegistration) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.requestCount["SaveValidatorRegistration"]++
	ds.entries[entry.Message.Pubkey] = &entry
	return nil
}

func (ds *MemoryDatastore) SaveValidatorRegistrations(entries []types.SignedValidatorRegistration) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.requestCount["SaveValidatorRegistrations"]++
	for _, entry := range entries {
		ds.entries[entry.Message.Pubkey] = &entry
	}
	return nil
}

// GetRequestCount returns the number of Request made to a method
func (ds *MemoryDatastore) GetRequestCount(method string) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return ds.requestCount[method]
}

func NewMemoryDatastore() *MemoryDatastore {
	return &MemoryDatastore{
		entries: make(map[types.PublicKey]*types.SignedValidatorRegistration),
		requestCount: make(map[string]int),
	}
}