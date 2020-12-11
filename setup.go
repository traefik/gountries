package gountries

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// New creates an Query object and unmarshal the json file.
func New() *Query {
	return NewFromPath(filepath.Join("data", "yaml"))
}

// NewFromPath creates a Query object from data folder in provided path.
func NewFromPath(dataPath string) *Query {
	queryInitOnce.Do(func() {
		queryInstance = &Query{
			Countries: populateCountries(dataPath),
		}

		queryInstance.NameToAlpha2 = populateNameIndex(queryInstance.Countries)
		queryInstance.Alpha3ToAlpha2 = populateAlphaIndex(queryInstance.Countries)
		queryInstance.NativeNameToAlpha2 = populateNativeNameIndex(queryInstance.Countries)

		subdivisions := populateSubdivisions(dataPath)

		for k := range queryInstance.Countries {
			c := queryInstance.Countries[k]

			c.subdivisions = subdivisions[strings.ToLower(c.Alpha2)]
			c.nameToSubdivision = map[string]SubDivision{}
			c.codeToSubdivision = map[string]SubDivision{}

			for _, s := range c.subdivisions {
				for _, n := range s.Names {
					c.nameToSubdivision[strings.ToLower(n)] = s
				}
				c.nameToSubdivision[strings.ToLower(s.Name)] = s
				c.codeToSubdivision[strings.ToLower(s.Code)] = s
			}

			queryInstance.Countries[k] = c
		}
	})

	return queryInstance
}

func populateNameIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)

	for alpha2, country := range countries {
		index[strings.ToLower(country.Name.Common)] = alpha2
		index[strings.ToLower(country.Name.Official)] = alpha2
	}

	return index
}

func populateAlphaIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)

	for alpha2, country := range countries {
		index[country.Codes.Alpha3] = alpha2
	}

	return index
}

func populateCountries(dataPath string) map[string]Country {
	// Try packed data first before custom data directory
	if yamlFileList, err := AssetDir("data/yaml/countries"); err == nil {
		return populateCountriesFromPackedData(yamlFileList, "data/yaml/countries")
	}

	countriesPath := path.Join(dataPath, "countries")

	info, err := ioutil.ReadDir(countriesPath)
	if err != nil {
		panic(fmt.Errorf("error loading Countries: %w", err))
	}

	countries := make(map[string]Country)

	for _, v := range info {
		if v.IsDir() {
			continue
		}

		var file []byte
		file, err = ioutil.ReadFile(filepath.Join(countriesPath, v.Name()))
		if err != nil {
			continue
		}

		country := Country{}
		err = yaml.Unmarshal(file, &country)
		if err != nil {
			continue
		}

		// Save
		countries[country.Codes.Alpha2] = country
	}

	return countries
}

func populateCountriesFromPackedData(fileList []string, path string) map[string]Country {
	countries := make(map[string]Country)

	for _, yamlFile := range fileList {
		data, err := Asset(filepath.Join(path, yamlFile))
		if err != nil {
			continue
		}

		var country Country
		if err = yaml.Unmarshal(data, &country); err != nil {
			continue
		}

		countries[country.Codes.Alpha2] = country
	}

	return countries
}

func populateSubdivisions(dataPath string) map[string][]SubDivision {
	// Try packed data first before custom data directory
	if yamlFileList, err := AssetDir("data/yaml/subdivisions"); err == nil {
		return populateSubdivisionsFromPackedData(yamlFileList, "data/yaml/subdivisions")
	}

	subdivisionsPath := path.Join(dataPath, "subdivisions")

	info, err := ioutil.ReadDir(subdivisionsPath)
	if err != nil {
		panic(fmt.Errorf("error loading Subdivisions: %w", err))
	}

	list := map[string][]SubDivision{}

	for _, v := range info {
		if v.IsDir() {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(subdivisionsPath, v.Name()))
		if err != nil {
			continue
		}

		var subdivisions []SubDivision
		err = yaml.Unmarshal(file, &subdivisions)
		if err != nil {
			continue
		}

		// Save
		// subdivisions = append(subdivisions, subdivision...)
		list[strings.Split(v.Name(), ".")[0]] = subdivisions
	}

	return list
}

func populateSubdivisionsFromPackedData(fileList []string, path string) map[string][]SubDivision {
	sd := make(map[string][]SubDivision)

	for _, yamlFile := range fileList {
		data, err := Asset(filepath.Join(path, yamlFile))
		if err != nil {
			continue
		}

		var subdivisions []SubDivision
		if err = yaml.Unmarshal(data, &subdivisions); err != nil {
			continue
		}

		alpha2 := strings.Split(yamlFile, ".")[0]

		sd[alpha2] = subdivisions
	}

	return sd
}

func populateNativeNameIndex(countries map[string]Country) map[string]string {
	index := make(map[string]string)

	for alpha2, country := range countries {
		for _, nativeNames := range country.Name.Native {
			index[strings.ToLower(nativeNames.Common)] = alpha2
			index[strings.ToLower(nativeNames.Official)] = alpha2
		}
	}

	return index
}
