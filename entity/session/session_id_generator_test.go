package session

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	fmt.Println(s)
	assert.True(t, len(s) == 32)
}
