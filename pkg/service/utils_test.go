package service

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeBasePathSpecialCaseWithSuccess(t *testing.T) {
	assert := assert.New(t)
	data := "hello world it's '/' test"
	expect := "hello world it's '-' test"

	data = filepath.Base(data)
	result := sanitizeBasePathSpecialCase(data)
	assert.Equal(expect, result)
}
