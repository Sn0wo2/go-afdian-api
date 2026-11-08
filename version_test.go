package afdian

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	ua := userAgent()

	assert.True(t, strings.HasPrefix(ua, "go-afdian-api/"), "unexpected user agent: %s", ua)

	assert.Equal(t, ua, userAgent())
}
