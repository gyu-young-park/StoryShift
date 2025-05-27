package servicevelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVelogUserExists(t *testing.T) {
	v := VelogService{}
	assert.Equal(t, true, v.IsVelogUserExists("chappi"))
	assert.Equal(t, false, v.IsVelogUserExists("qowiendsm192j3i1"))
}
