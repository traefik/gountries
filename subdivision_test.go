package gountries

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubdivisions(t *testing.T) {
	se, err := query.FindCountryByAlpha("SWE")
	require.NoError(t, err)

	subDiv := se.SubDivisions()

	assert.Len(t, subDiv, 21)

	found, err := se.FindSubdivisionByName(subDiv[0].Name)
	require.NoError(t, err)
	assert.Equal(t, subDiv[0], found)

	for _, n := range subDiv[0].Names {
		found, err = se.FindSubdivisionByName(n)
		require.NoError(t, err)
		assert.Equal(t, subDiv[0], found)
	}

	found, err = se.FindSubdivisionByCode(subDiv[0].Code)
	require.NoError(t, err)
	assert.Equal(t, subDiv[0], found)
}
