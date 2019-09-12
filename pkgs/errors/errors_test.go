package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	c "github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
)

func TestError(t *testing.T) {
	err := ErrNoPermissions
	errStrign := err.Error()
	assert.Equal(t, c.NoPermissions, errStrign)
}
