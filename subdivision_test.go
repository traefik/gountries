package gountries

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubdivisions(t *testing.T) {

	se, err := query.FindCountryByAlpha("SWE")
	require.NoError(t, err)

	subd := se.SubDivisions()

	assert.Len(t, subd, 21)

	found, err := se.FindSubdivisionByName(subd[0].Name)
	require.NoError(t, err)
	assert.Equal(t, subd[0], found)

	for _, n := range subd[0].Names {
		found, err = se.FindSubdivisionByName(n)
		require.NoError(t, err)
		assert.Equal(t, subd[0], found)
	}

	found, err = se.FindSubdivisionByCode(subd[0].Code)
	require.NoError(t, err)
	assert.Equal(t, subd[0], found)
}
