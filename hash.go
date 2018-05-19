package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/dchest/siphash"
)

// HashTracker tracks the input and the output for a hash function
type HashTracker struct {
	Geos          []GeoPoint
	Zips          []ZipCode
	Words         []Word
	rawHashValues hashValues
}

type hashValues struct {
	Geos  map[string]string `json:"geo"`
	Zips  map[int]string    `json:"zip"`
	Words map[string]string `json:"word"`
}

type hashFunc func(key interface{}) (hash string)

var allHashes = map[string]hashFunc{
	"sha256":  hashSha256,
	"sha1":    hashSha1,
	"siphash": hashSipHash,
	"md5":     hashMd5,
}

// NewHashTracker builds all the values that will be hashed
func NewHashTracker(geos []GeoPoint, zips []ZipCode, words []Word) *HashTracker {
	return &HashTracker{
		Geos:  geos,
		Zips:  zips,
		Words: words,
		rawHashValues: hashValues{
			Geos:  map[string]string{},
			Zips:  map[int]string{},
			Words: map[string]string{},
		},
	}
}

// Hash hashes all the values in the hash tracker and records their values
func (ht *HashTracker) Hash(hashName string) *HashTracker {
	for _, geo := range ht.Geos {
		hashVal := geo.Hash(hashName)
		gBytes, err := json.Marshal(geo)
		if err != nil {
			panic(err)
		}

		ht.rawHashValues.Geos[string(gBytes)] = hashVal
	}

	for _, zip := range ht.Zips {
		hashVal := zip.Hash(hashName)
		ht.rawHashValues.Zips[int(zip)] = hashVal
	}

	for _, word := range ht.Words {
		hashVal := word.Hash(hashName)
		ht.rawHashValues.Words[string(word)] = hashVal
	}
	return ht
}

// Export exports all the raw hashes and their mapping
func (ht *HashTracker) Export() (rawOutput []byte, err error) {
	return json.Marshal(ht.rawHashValues)
}

func hashSha256(key interface{}) (hash string) {
	switch k := key.(type) {
	case int:
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(k))
		return fmt.Sprintf("%x", sha256.Sum256(bs))
	case []byte:
		return fmt.Sprintf("%x", sha256.Sum256(k))
	case string:
		return fmt.Sprintf("%x", sha256.Sum256([]byte(k)))
	}

	panic(fmt.Sprintf("unsupported hash type: %T", key))
}

func hashSha1(key interface{}) (hash string) {
	switch k := key.(type) {
	case int:
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(k))
		return fmt.Sprintf("%x", sha1.Sum(bs))
	case []byte:
		return fmt.Sprintf("%x", sha1.Sum(k))
	case string:
		return fmt.Sprintf("%x", sha1.Sum([]byte(k)))
	}

	panic(fmt.Sprintf("unsupported hash type: %T", key))
}

func hashMd5(key interface{}) (hash string) {
	switch k := key.(type) {
	case int:
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(k))
		return fmt.Sprintf("%x", md5.Sum(bs))
	case []byte:
		return fmt.Sprintf("%x", md5.Sum(k))
	case string:
		return fmt.Sprintf("%x", md5.Sum([]byte(k)))
	}

	panic(fmt.Sprintf("unsupported hash type: %T", key))
}

func hashSipHash(key interface{}) (hash string) {
	switch k := key.(type) {
	case int:
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(k))
		return fmt.Sprintf("%x", siphash.Hash(0, 2048, bs))
	case []byte:
		return fmt.Sprintf("%x", siphash.Hash(0, 2048, k))
	case string:
		return fmt.Sprintf("%x", siphash.Hash(0, 2048, []byte(k)))
	}

	panic(fmt.Sprintf("unsupported hash type: %T", key))
}
