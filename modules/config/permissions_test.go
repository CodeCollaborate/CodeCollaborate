package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissionLabel(t *testing.T) {
	for label, level := range PermissionsByLabel {
		permission, err := PermissionByLevel(level)
		assert.Nil(t, err)
		assert.Equal(t, label, permission.Label, "unexpected permission label")
	}
}

func TestGetPermissionLevel(t *testing.T) {
	for label, level := range PermissionsByLabel {
		permission, err := PermissionByLabel(label)
		assert.Nil(t, err)
		assert.Equal(t, level, permission.Level, "unexpected permission label")
	}
}
