package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Opts are the options for the script
type Opts struct {
	HashName   string `long:"hash" description:"the hash function you want to use to hash all the values" required:"true"`
	OutputPath string `short:"o" description:"the output path to write the hashes to" default:"output.json"`
}

var opts Opts
var parser = flags.NewParser(&opts, flags.Default)

func main() {
	_, err := parser.Parse()
	if err != nil {
		if isUsage(err) || isCommand(err) {
			os.Exit(1)
		}
		log.Fatalf("Error parsing arguments: %s", err)
	}

	if _, found := allHashes[opts.HashName]; !found {
		hashes := []string{}
		for hash := range allHashes {
			hashes = append(hashes, hash)
		}
		log.Fatalf("invalid hash given. Possible values are %s", strings.Join(hashes, ", "))
	}

	lines, err := readNewlineFile("data/words.txt")
	if err != nil {
		log.Fatalf("Error reading newline file: %v\n", err)
	}

	words, err := ParseWords(lines)
	if err != nil {
		log.Fatalf("Error parsing words out of newline file: %v\n", err)
	}

	csvLines, err := readCSVFile("data/zipcodes.txt")
	if err != nil {
		log.Fatalf("Error reading csv file: %v\n", err)
	}

	geos, err := ParseGeoPoints(csvLines[1:])
	if err != nil {
		log.Fatalf("Error parsing geos out of csv file: %v\n", err)
	}

	zips, err := ParseZipCodes(csvLines[1:])
	if err != nil {
		log.Fatalf("Error parsing zips out of csv file: %v\n", err)
	}

	exportedHashes, err := NewHashTracker(geos, zips, words).Hash(opts.HashName).Export()
	if err != nil {
		log.Fatalf("Error creating and exporting hashes: %v\n", err)
	}

	ioutil.WriteFile(opts.OutputPath, exportedHashes, 0666)
}

func readNewlineFile(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if s.Err() != nil {
		return nil, fmt.Errorf("Unable to scan newline file: %v", err)
	}

	return lines, nil

}

func readCSVFile(path string) (lines [][]string, err error) {
	f, err := os.Open("data/zipcodes.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening csv file: %v", err)
	}
	defer f.Close()

	lines, err = csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv file: %v", err)
	}

	return lines, nil
}

func isUsage(err error) bool {
	return strings.HasPrefix(err.Error(), "Usage:")
}

func isCommand(err error) bool {
	return strings.HasPrefix(err.Error(), "Please specify")
}
