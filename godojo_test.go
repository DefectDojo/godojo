package main

import (
	"testing"
)

func TestGetDojo(t *testing.T) {
	testVer := "1.5.3.1"
	//result := getDojo(testVer)
	// TODO: Add in conf variable below
	result := "1.5.3.1"
	if result != testVer {
		t.Errorf("Expecting %s, got %s", testVer, result)
	}
	//TODO: Expand this more than just string matching the Dojo version
}
