package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ablease/grn/tilefetcher"
	"github.com/ablease/grn/unzip"
)

func main() {
	tileVersion := flag.String("tile-version", "1.18", "Tile Version you want to generate release notes for.")

	flag.Parse()

	fmt.Println(GenerateReleaseNotes())
	fmt.Println("tile-version:", *tileVersion)

	// create working directory
	path := "/tmp/grn/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	// cleanup
	// defer rmDir(path)

	// Fetch tile
	tile, err := tilefetcher.Fetch(path)
	if err != nil {
		log.Fatal(err)
	}

	//Extract releases and release versions from tile tile = archive, path = destination
	files, err := unzip.Do(tile, path)
	if err != nil {
		log.Fatal(err)
	}

	releaseFiles := releaseFiles(files)

	releaseVersions := extractReleaseVersions(releaseFiles)

	for release, version := range releaseVersions {
		fmt.Printf("Release Name: %s, Version: %s \n", release, version)
	}
}

func GenerateReleaseNotes() string {
	return "Generate Release Notes"
}

func rmDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.RemoveAll(path)
	}
}

func releaseFiles(files []string) []string {
	// go has no generics so we need to extract releases from the list of files
	releases := make([]string, 0)
	for _, file := range files {
		if isRelease(file) {
			releases = append(releases, file)
		}
	}
	return releases
}

func isRelease(file string) bool {
	return strings.Contains(file, "releases/release")
}

func extractReleaseVersions(files []string) map[string]string {
	// We want release name + release version
	rvs := make(map[string]string)

	for _, file := range files {
		release := strings.TrimPrefix(file, "/tmp/grn/releases/release-")

		versionRegex := regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)
		version := versionRegex.Find([]byte(release))
		rvs[release] = string(version)
	}

	return rvs
}

// 1. Download a built tile
//      Specify a tile minor (1.15, 1.16, 1.17)
//      Tiles are uploaded to s3
//      Download the latest minor that the passed the product pipeline
// 2. Dissect tile
// 3. Extract release names and versions
