package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissionLabel(t *testing.T) {
	for _, perm := range permissions {
		label, err := GetPermissionLabel(perm.Level)
		assert.Nil(t, err)
		assert.Equal(t, perm.Label, label, "unexpected permission label")
	}
}

func TestGetPermissionLevel(t *testing.T) {
	for _, perm := range permissions {
		level, err := GetPermissionLevel(perm.Label)
		assert.Nil(t, err)
		assert.Equal(t, perm.Level, level, "unexpected permission label")
	}
}
