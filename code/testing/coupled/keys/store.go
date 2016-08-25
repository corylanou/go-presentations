package keys

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// START STORE-OMIT
type Store struct {
	id     int
	values map[string]interface{}
	mu     sync.RWMutex
}

// END STORE-OMIT

func NewStore() *Store {
	s := Store{}
	s.values = make(map[string]interface{})
	return &s
}

// START STORE-SET-OMIT
func (vs *Store) Set(key string, value interface{}) {
	vs.id++
	// Make this an asynchronous call
	go func() {
		vs.mu.Lock()
		defer vs.mu.Unlock()

		// take a random amout of time to return to simulate the real world
		rand.Seed(time.Now().Unix())
		t := time.Duration(rand.Intn(3)) * time.Second
		time.Sleep(t)

		vs.values[key] = value
		log.Println("inserted: ", key, " with value of ", value)
	}()
}

// END STORE-SET-OMIT

// START STORE-GET-OMIT
func (vs *Store) Get(key string) (interface{}, error) {
	if v, ok := vs.values[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("%s not found", key)
}

// END STORE-GET-OMIT

func (vs *Store) Count() int {
	return len(vs.values)
}
