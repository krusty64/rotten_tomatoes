package rotten_tomatoes

import (
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	c, err := InitSetup()
	if err != nil {
		t.Error(err)
	}

	if c.Links.Movies == "" {
		t.Error("setup not initialized:", c.Links.Movies)
	}

	if c.LinkTemplate == "" {
		t.Error("Link template not initialized", c.LinkTemplate)
	}

	q := c.LinkUrl.Query()
	for k, v := range q {
		if len(v) != 1 {
			t.Error("Invalid query:", k, v)
		}
		if strings.HasPrefix(v[0], "{") && k != "q" {
			t.Error("missing default value", k, v)
		}
	}
}
