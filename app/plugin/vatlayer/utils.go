package vatlayer

type VatlayerConfiguration struct {
	AccessKey           string
	ExcludedCountries   []string
	CountriesFromOrigin []string
	OriginCountry       string
}
