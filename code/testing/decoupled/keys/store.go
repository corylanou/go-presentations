package keys

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Store struct {
	id     int
	values map[string]interface{}
	mu     sync.RWMutex
}

func NewStore() *Store {
	s := Store{}
	s.values = make(map[string]interface{})
	return &s
}

func (vs *Store) Set(key string, value interface{}) {
	vs.id++
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

func (vs *Store) Get(key string) (interface{}, error) {
	if v, ok := vs.values[key]; ok {
		return v, nil
	}
	return nil, notFound{key: key}
}

func (vs *Store) Count() int {
	return len(vs.values)
}
