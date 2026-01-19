package domain

import (
	"encoding/json"
	"fmt"
)

var (
	ADMIN     Role = Role{ID: 1, Name: "admin"}
	SUPPORTER Role = Role{ID: 2, Name: "supporter"}
	CREATOR   Role = Role{ID: 3, Name: "creator"}
)

type Role struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Roles []Role

func (r *Roles) Scan(value any) error {
	if value == nil {
		*r = make([]Role, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for JSONB: %T", value)
	}
	return json.Unmarshal(bytes, r)
}
