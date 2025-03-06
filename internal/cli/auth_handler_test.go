package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMissingAuth(t *testing.T) {
	err := handleMissingAuth()

	assert.Error(t, err)
	assert.Equal(t, "authentication required", err.Error())
}
