package data

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"sync"
)

var once sync.Once

//TODO I have some question around best practices here, from the testing I've done I don't think sync.Mutex or sync.Map
//is required since we're using the singelton pattern.
type State struct {
	hashCount     int64
	hashCountLock sync.Mutex
	// This might not be the right data structure for this.
	hashes sync.Map
}

var instance *State

func Get() *State {
	once.Do(func() {
		instance = &State{
			hashCount:     0,
			hashCountLock: sync.Mutex{},
			hashes:        sync.Map{},
		}
	})
	return instance
}

// GetNextIdentifier Manages locking/unlocking the auto incremented identifier
// The identifier should be refactored to return a UUID to help avoid ID collisions
func (state *State) GetNextIdentifier() int64 {
	state.hashCountLock.Lock()
	state.hashCount++
	state.hashCountLock.Unlock()
	return state.hashCount
}

func (state *State) GetIdentifierCount() int64 {
	//TODO TW: Should we use RWMutex for this?
	return state.hashCount
}

// GetHashedPassword This returns the value saved in the map if found = true otherwise it'll return an empty string.
func (state *State) GetHashedPassword(identifier int64) (value string, found bool) {
	result, ok := state.hashes.Load(identifier)
	return fmt.Sprintf("%s", result), ok
}

// SavePassword This will hash the password and store it in a map
func (state *State) SavePassword(identifier int64, password string) {
	value := hash(password)
	state.hashes.Store(identifier, base64Encode(value))
}

func base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func hash(value string) string {
	h := sha512.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}
