package config

import "errors"

/**
 * Permission constants for API access levels
 */
var permissions = []Permission{
	{Label: "Read", Level: 1},
	{Label: "Write", Level: 4},
	{Label: "Admin", Level: 8},
	{Label: "Owner", Level: 10},
}

// Permission is the struct representation of an API permission level
type Permission struct {
	Level int8
	Label string
}

// internal storage maps
var byLabel map[string]Permission
var byLevel map[int8]Permission

// initialize maps
func init() {
	byLabel = make(map[string]Permission)
	byLevel = make(map[int8]Permission)

	for _, perm := range permissions {
		byLabel[perm.Label] = perm
		byLevel[perm.Level] = perm
	}
}

// ErrNoMatchingPermission is returned if a permission that does not exist is attempted to be accessed
var ErrNoMatchingPermission = errors.New("Not a valid server permission level")

// PermissionByLevel returns the string representation of the provided level, if found
func PermissionByLevel(level int8) (Permission, error) {
	label, ok := byLevel[level]
	if !ok {
		return Permission{}, ErrNoMatchingPermission
	}
	return label, nil
}

// PermissionByLabel returns the int8 representation of the provided label, if found
func PermissionByLabel(label string) (Permission, error) {
	level, ok := byLabel[label]
	if !ok {
		return Permission{}, ErrNoMatchingPermission
	}
	return level, nil
}
