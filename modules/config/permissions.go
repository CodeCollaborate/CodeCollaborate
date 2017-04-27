package config

import "errors"

// PermissionsByLabel is the permission constants for API access levels
var PermissionsByLabel = map[string]int{
	"read":  1,
	"write": 4,
	"admin": 8,
	"owner": 10,
}

// Permission is the struct representation of an API permission level
type Permission struct {
	Level int
	Label string
}

// internal map in other direction
var byLevel map[int]string

// initialize maps
func init() {
	byLevel = make(map[int]string)
	for label, level := range PermissionsByLabel {
		byLevel[level] = label
	}
}

// ErrNoMatchingPermission is returned if a permission that does not exist is attempted to be accessed
var ErrNoMatchingPermission = errors.New("Not a valid server permission level")

// PermissionByLevel returns the string representation of the provided level, if found
func PermissionByLevel(level int) (Permission, error) {
	label, ok := byLevel[level]
	if !ok {
		return Permission{}, ErrNoMatchingPermission
	}
	return Permission{
		Label: label,
		Level: level,
	}, nil
}

// PermissionByLabel returns the int representation of the provided label, if found
func PermissionByLabel(label string) (Permission, error) {
	level, ok := PermissionsByLabel[label]
	if !ok {
		return Permission{}, ErrNoMatchingPermission
	}
	return Permission{
		Level: level,
		Label: label,
	}, nil
}
