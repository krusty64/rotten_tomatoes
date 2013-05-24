package rotten_tomatoes

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitSetup()
	t.Error(url_setup)
}
