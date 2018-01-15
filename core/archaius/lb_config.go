package archaius

import "strings"

// GetServerListFilters get server list filters
func GetServerListFilters() []string {
	return strings.Split(GetString(GetFilterNamesKey(), ""), ",")
}
