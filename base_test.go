package myminio

import (
	"testing"
)

func Test_MM(t *testing.T) {
	c, err := NewClient("", "", "", "")
	t.Log("c,err", c, err)
}
