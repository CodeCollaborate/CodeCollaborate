package config

import "errors"

/**
 * Permission constants for API access levels
 */
var permissions = []permission{
	{Label: "Read", Level: 1},
	{Label: "Write", Level: 4},
	{Label: "Admin", Level: 8},
	{Label: "Owner", Level: 10},
}

// permission is the struct representation of an API permission level
type permission struct {
	Level int8
	Label string
}

// internal storage maps
var byLabel map[string]int8
var byLevel map[int8]string

// initialize maps
func init() {
	byLabel = make(map[string]int8)
	byLevel = make(map[int8]string)

	for _, perm := range permissions {
		byLabel[perm.Label] = perm.Level
		byLevel[perm.Level] = perm.Label
	}
}

// ErrNoPermission is returned if a permission that does not exist is attempted to be accessed
var ErrNoPermission = errors.New("Not a valid server permission level")

// GetPermissionLabel returns the string representation of the provided level, if found
func GetPermissionLabel(level int8) (string, error) {
	label, ok := byLevel[level]
	if !ok {
		return "", ErrNoPermission
	}
	return label, nil
}

// GetPermissionLevel returns the int8 representation of the provided label, if found
func GetPermissionLevel(label string) (int8, error) {
	level, ok := byLabel[label]
	if !ok {
		return -1, ErrNoPermission
	}
	return level, nil
}
