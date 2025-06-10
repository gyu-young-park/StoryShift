package servicevelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeBasePathSpecialCaseWithSuccess(t *testing.T) {
	data := "hello world it's '/' test"
	expect := "hello world it's '-' test"

	sanitized, isSanitized := sanitizeBasePathSpecialCase(data)

	assert.Equal(t, expect, sanitized)
	assert.True(t, isSanitized)
}
