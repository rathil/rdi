package standard

import (
	"sync"

	"github.com/rathil/rdi"
)

type di struct {
	storage sync.Map
	parent  rdi.DI
}
