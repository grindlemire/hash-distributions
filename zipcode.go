package main

import "strconv"

// ZipCode represents a hashable zipcode
type ZipCode int

// ParseZipCodes parses the zip codes out of the csv reader
func ParseZipCodes(lines [][]string) (zips []ZipCode, err error) {
	for _, line := range lines {
		zip, err := strconv.Atoi(line[0])
		if err != nil {
			return zips, err
		}
		zips = append(zips, ZipCode(zip))
	}
	return zips, nil
}

// Hash hashes the zipcode using the specified hash
func (z ZipCode) Hash(hashName string) (hashVal string) {
	return allHashes[hashName](int(z))
}
