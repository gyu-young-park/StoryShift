package servicevelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVelogUserExists(t *testing.T) {
	assert.Equal(t, true, IsVelogUserExists("chappi"))
	assert.Equal(t, false, IsVelogUserExists("qowiendsm192j3i1"))
}
