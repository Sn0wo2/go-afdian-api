package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	ua := UserAgent{}.String()

	assert.True(t, strings.HasPrefix(ua, "go-afdian-api/"), "unexpected user agent: %s", ua)

	assert.Equal(t, ua, UserAgent{}.String())
}
