package velog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSeriesAPIWhenSuccess(t *testing.T) {
	velogApi := NewVelogAPI("https://v2.velog.io/graphql")
	readSeries, err := velogApi.ReadSeries("chappi", "CKA")

	assert.NoError(t, err, "error occured")
	t.Log(readSeries)
}
