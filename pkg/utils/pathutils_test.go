package utils

import "testing"

func TestPathParser(t *testing.T) {
	path := "conn/dbname/schemaname/tabname"
	parts := ParsePath(path)

	if parts.ConfigConnection != "conn" {
		t.Fail()
	}
	if parts.DatabaseName != "dbname" {
		t.Fail()
	}
	if parts.SchemaName != "schemaname" {
		t.Fail()
	}
	if parts.TableName != "tabname" {
		t.Fail()
	}
}
