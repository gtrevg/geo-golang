// Package nominatim is a geo-golang based MapRequest Nominatim geocode/reverse geocode client
package nominatim

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/osm"
)

type (
	baseURL string

	geocodeResponse struct {
		DisplayName     string `json:"display_name"`
		Lat, Lon, Error string
		Addr            osm.Address `json:"address"`
	}
)

var key string

// Geocoder constructs MapRequest Nominatim geocoder
func Geocoder(k string, baseURLs ...string) geo.Geocoder {
	key = k
	return geo.HTTPGeocoder{
		EndpointBuilder:       baseURL(getUrl(baseURLs...)),
		ResponseParserFactory: func() geo.ResponseParser { return &geocodeResponse{} },
	}
}

func getUrl(baseURLs ...string) string {
	if len(baseURLs) > 0 {
		return baseURLs[0]
	}
	return "http://open.mapquestapi.com/nominatim/v1/"
}

func (b baseURL) GeocodeURL(address string) string {
	return string(b) + "search.php?key=" + key + "&format=json&limit=1&q=" + address
}

func (b baseURL) ReverseGeocodeURL(l geo.Location) string {
	return string(b) + "reverse.php?key=" + key + fmt.Sprintf("&format=json&lat=%f&lon=%f", l.Lat, l.Lng)
}

func (r *geocodeResponse) Location() (*geo.Location, error) {
	if r.Error != "" {
		return nil, fmt.Errorf("geocode error: %s", r.Error)
	}

	return &geo.Location{
		Lat: parseFloat(r.Lat),
		Lng: parseFloat(r.Lon),
	}, nil
}

func (r *geocodeResponse) Address() (*geo.Address, error) {
	if r.Error != "" {
		return nil, fmt.Errorf("reverse geocode error: %s", r.Error)
	}

	return &geo.Address{
		FormattedAddress: r.DisplayName,
		HouseNumber:      r.Addr.HouseNumber,
		Street:           r.Addr.Street(),
		Suburb:           r.Addr.Suburb,
		City:             r.Addr.Locality(),
		State:            r.Addr.State,
		County:           r.Addr.County,
		Postcode:         r.Addr.Postcode,
		Country:          r.Addr.Country,
		CountryCode:      strings.ToUpper(r.Addr.CountryCode),
	}, nil
}

func parseFloat(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}
