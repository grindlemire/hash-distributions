package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

// GeoPoint represents a hashable geopoint
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// ParseGeoPoints parses the geopoints out of the csv reader
func ParseGeoPoints(lines [][]string) (geos []GeoPoint, err error) {
	for _, line := range lines {

		lat, err := strconv.ParseFloat(strings.TrimSpace(line[1]), 64)
		if err != nil {
			return geos, err
		}

		lon, err := strconv.ParseFloat(strings.TrimSpace(line[2]), 64)
		if err != nil {
			return geos, err
		}

		geos = append(geos, GeoPoint{Lat: lat, Lon: lon})
	}
	return geos, nil
}

// Hash hashes the GeoPoint using the specified hash
func (g GeoPoint) Hash(hashName string) (hashVal string) {
	gBytes, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	return allHashes[hashName](gBytes)
}
