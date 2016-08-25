package keys

import "fmt"

type notFound struct {
	key string
}

func (nf notFound) NotFound() {}

func (nf notFound) Error() string {
	return fmt.Sprintf("%q not found", nf.key)
}
