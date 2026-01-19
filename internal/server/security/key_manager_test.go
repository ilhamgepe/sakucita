package security

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRSA(t *testing.T) {
	err := testSecurity.LoadRSAKeys("../../../keys")
	if err != nil {
		assert.NoError(t, err)
	}

	for _, key := range testSecurity.rsaKeys {
		fmt.Printf("private key %+v\n\n\n", key.private)
		fmt.Printf("public key %+v\n\n\n", key.public)
	}
	fmt.Printf("active kid %+v\n", testSecurity.activeKID)
}
