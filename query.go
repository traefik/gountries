package gountries

import (
	"fmt"
	"strings"
	"sync"
)

// Query holds a reference to the QueryHolder struct.
var (
	queryInitOnce sync.Once
	queryInstance *Query
)

// Query contains queries for countries, cities, etc.
type Query struct {
	Countries          map[string]Country
	NameToAlpha2       map[string]string
	Alpha3ToAlpha2     map[string]string
	NativeNameToAlpha2 map[string]string
}

// FindCountryByName finds a country by given name.
func (q *Query) FindCountryByName(name string) (result Country, err error) {
	alpha2, exists := q.NameToAlpha2[strings.ToLower(name)]
	if !exists {
		return Country{}, fmt.Errorf("could not find country with name: %s", name)
	}

	return q.Countries[alpha2], nil
}

// FindCountryByNativeName finds a country by given native name.
func (q *Query) FindCountryByNativeName(name string) (result Country, err error) {
	alpha2, exists := q.NativeNameToAlpha2[strings.ToLower(name)]
	if !exists {
		return Country{}, fmt.Errorf("could not find country with native name: %s", name)
	}

	return q.Countries[alpha2], nil
}

// FindCountryByAlpha finds a country by given code.
func (q *Query) FindCountryByAlpha(code string) (Country, error) {
	codeU := strings.ToUpper(code)

	switch len(code) {
	case 2:
		country, exists := q.Countries[codeU]
		if !exists {
			return Country{}, fmt.Errorf("could not find country with code %s", code)
		}
		return country, nil

	case 3:
		alpha2, exists := q.Alpha3ToAlpha2[codeU]
		if !exists {
			return Country{}, fmt.Errorf("could not find country with code: %s", code)
		}
		return q.Countries[alpha2], nil

	default:
		return Country{}, fmt.Errorf("invalid code format: %s", code)
	}
}

// FindAllCountries returns a list of all countries.
func (q *Query) FindAllCountries() map[string]Country {
	return q.Countries
}

// FindCountries finds a Country based on the given struct data.
func (q Query) FindCountries(c Country) []Country {
	var countries []Country

	for _, country := range q.Countries {
		// Name
		if c.Name.Common != "" && strings.EqualFold(c.Name.Common, country.Name.Common) {
			continue
		}

		// Alpha
		if c.Alpha2 != "" && c.Alpha2 != country.Alpha2 {
			continue
		}
		if c.Alpha3 != "" && c.Alpha3 != country.Alpha3 {
			continue
		}

		// Geo
		if c.Geo.Continent != "" && !strings.EqualFold(c.Geo.Continent, country.Geo.Continent) {
			continue
		}
		if c.Geo.Region != "" && !strings.EqualFold(c.Geo.Region, country.Geo.Region) {
			continue
		}
		if c.Geo.SubRegion != "" && !strings.EqualFold(c.Geo.SubRegion, country.Geo.SubRegion) {
			continue
		}

		// Misc
		if c.InternationalPrefix != "" && !strings.EqualFold(c.InternationalPrefix, country.InternationalPrefix) {
			continue
		}

		// Bordering countries

		allMatch := false

		if len(c.BorderingCountries()) > 0 {
			for _, c1 := range c.BorderingCountries() {
				match := false

				for _, c2 := range country.BorderingCountries() {
					match = c1.Alpha2 == c2.Alpha2
					if match {
						break
					}
				}

				if match {
					allMatch = true
				} else {
					allMatch = false
					break
				}
			}

			if !allMatch {
				continue
			}
		}

		// Append if all matches
		countries = append(countries, country)
	}

	return countries
}
