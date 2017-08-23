package objects

import (
	"encoding/json"
	"fmt"
)

// Indirect -
type Indirect struct {
	ObjectNumber     uint
	GenerationNumber uint
}

func (s Indirect) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d %d R", s.ObjectNumber, s.GenerationNumber))
}
