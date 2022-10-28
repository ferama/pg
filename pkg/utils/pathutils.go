package utils

import "strings"

type PathParts struct {
	ConfigConnection string
	DatabaseName     string
	SchemaName       string
	TableName        string
}

func ParsePath(path string, requireTrailingSlash bool) *PathParts {
	parts := strings.Split(path, "/")
	if requireTrailingSlash && !strings.HasSuffix(path, "/") {
		parts = parts[:len(parts)-1]
	}
	pp := &PathParts{}
	if len(parts) > 0 {
		pp.ConfigConnection = parts[0]
	}
	if len(parts) > 1 {
		pp.DatabaseName = parts[1]
	}
	if len(parts) > 2 {
		pp.SchemaName = parts[2]
	}
	if len(parts) > 3 {
		pp.TableName = parts[3]
	}

	return pp
}
