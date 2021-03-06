package gountries

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCountryByName(t *testing.T) {
	// Test for lowercase
	//

	result, err := query.FindCountryByName("sweden")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Lowercase country names should match")

	// Test for common name
	result, err = query.FindCountryByName("United States")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "US", "Lowercase country names should match")

	// Test for official name
	result, err = query.FindCountryByName("United States of America")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "US", "Lowercase country names should match")

	// Test for uppercase
	//

	result, err = query.FindCountryByName("SWEDEN")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Uppercase country names should match")

	// Test for invariants
	//

	invariants := []string{"Sweden", "SwEdEn", "SWEden"}

	for _, invariant := range invariants {
		result, err = query.FindCountryByName(invariant)
		require.NoError(t, err)

		assert.Equal(t, result.Alpha2, "SE", fmt.Sprintf("Invariants of country names, eg sWeden,SWEDEN,swEdEn should match. %s did not match", invariant))
	}
}

func TestFindCountryByAlpha(t *testing.T) {
	// Test for lowercase

	result, err := query.FindCountryByAlpha("se")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Lowercase country names should match")

	// Test for uppercase

	result, err = query.FindCountryByAlpha("SE")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Uppercase country names should match")

	// Test for invariants

	result, err = query.FindCountryByAlpha("Se")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Invariants of country names, eg sWeden,SWEDEN,swEdEn should match")

	// Test for wrong code types (wrong length)

	result, err = query.FindCountryByAlpha("SEE")
	require.Error(t, err)

	assert.EqualError(t, err, "could not find country with code: SEE")

	// Test for wrong code types: too long

	result, err = query.FindCountryByAlpha("SEEE")
	require.Error(t, err)

	assert.EqualError(t, err, "invalid code format: SEEE")

	// Test for wrong code types: too short

	result, err = query.FindCountryByAlpha("S")
	require.Error(t, err)
	assert.EqualError(t, err, "invalid code format: S")
}

func TestFindAllCountries(t *testing.T) {
	assert.Len(t, query.FindAllCountries(), 249)
}

func TestFindCountries(t *testing.T) {
	country := Country{}
	country.Alpha2 = "SE"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 1)

	assert.Equal(t, countries[0].Alpha2, "SE", fmt.Sprintf("Countries did not return expected result %s: %s", "SE", countries[0].Alpha2))
}

func TestFindCountriesByRegion(t *testing.T) {
	country := Country{}
	country.Geo.Region = "Europe"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 52) // 52 is not the exact number of countries in Europe. Fix this later
}

func TestFindCountriesByContinent(t *testing.T) {
	country := Country{}
	country.Geo.Continent = "Europe"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 52) // 52 is not the exact number of countries in Europe. Fix this later
}

func TestFindCountriesBySubRegion(t *testing.T) {
	country := Country{}
	country.Geo.SubRegion = "Eastern Asia"

	countries := query.FindCountries(country)

	assert.Len(t, countries, 8) // 8 is not the exact number of countries in Eastern Asia. Fix this later
}

func TestFindCountryByNativeName(t *testing.T) {
	// Test for common name
	//

	result, err := query.FindCountryByNativeName("Sverige")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Common native country names should match")

	// Test for common name
	result, err = query.FindCountryByNativeName("Konungariket Sverige")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Official native country names should match")

	// Test for lowercase
	//

	result, err = query.FindCountryByNativeName("sverige")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Uppercase native country names should match")

	// Test for uppercase
	//

	result, err = query.FindCountryByNativeName("SVERIGE")
	require.NoError(t, err)

	assert.Equal(t, result.Alpha2, "SE", "Uppercase native country names should match")

	// Test for invariants
	//

	invariants := []string{"sVEriGE", "SveRIge", "SVErige"}

	for _, invariant := range invariants {
		result, err = query.FindCountryByNativeName(invariant)
		require.NoError(t, err)

		assert.Equal(t, result.Alpha2, "SE", fmt.Sprintf("Invariants of native country names, eg sVEriGE,SveRIge,SVErige should match. %s did not match", invariant))
	}
}

func ExampleQuery_FindAllCountries_borderingCountries() {
	country := Country{
		Borders: []string{"DEU"},
	}

	countries := query.FindCountries(country)

	var c []string
	for _, country := range countries {
		c = append(c, country.Name.Common)
	}

	sort.Strings(c)

	for _, name := range c {
		fmt.Println(name)
	}

	// Output:
	// Austria
	// Belgium
	// Czech Republic
	// Denmark
	// France
	// Luxembourg
	// Netherlands
	// Poland
	// Switzerland
}

func ExampleQuery_FindAllCountries_borderingCountries2() {
	country := Country{
		Borders: []string{
			"DEU",
			"CHE",
		},
	}

	countries := query.FindCountries(country)
	var c []string
	for _, country := range countries {
		c = append(c, country.Name.Common)
	}
	sort.Strings(c)
	for _, name := range c {
		fmt.Println(name)
	}

	// Output:
	// Austria
	// France
}

var result Country

func BenchmarkCountryLookupByName(b *testing.B) {
	q := New()

	var names []string
	for key := range q.Countries {
		names = append(names, q.Countries[key].Name.Common)
	}

	for n := 0; n <= b.N; n++ {
		randIndex := rand.Intn(len(q.Countries))

		c, err := q.FindCountryByName(names[randIndex])
		if err != nil {
			b.Fail()
		}

		result = c
	}
}
