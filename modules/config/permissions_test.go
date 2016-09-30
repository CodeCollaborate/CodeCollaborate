package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissionLabel(t *testing.T) {
	for _, perm := range permissions {
		permission, err := PermissionByLevel(perm.Level)
		assert.Nil(t, err)
		assert.Equal(t, perm.Label, permission.Label, "unexpected permission label")
	}
}

func TestGetPermissionLevel(t *testing.T) {
	for _, perm := range permissions {
		permission, err := PermissionByLabel(perm.Label)
		assert.Nil(t, err)
		assert.Equal(t, perm.Level, permission.Level, "unexpected permission label")
	}
}
