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
	defer rmDir(path)

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

	// Output as HTML
	// <tr>
	// <td>on-demand-service-broker</td>
	// <td>0.36.0</td>
	// </tr>

	for release, version := range releaseVersions {
		fmt.Println("<tr>")
		fmt.Printf("    <td>%s</td>\n", release)
		fmt.Printf("    <td>%s</td>\n", version)
		fmt.Println("</tr>")
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
	// E.G turn  "/tmp/grn/releases/release-service-metrics-1.12.1.on-ubuntu-xenial-stemcell.315.99.tgz"
	// into "service-metrics 1.12.1"

	rvs := make(map[string]string)

	for _, file := range files {
		release := strings.TrimPrefix(file, "/tmp/grn/releases/release-")

		versionRegex := regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)
		v := versionRegex.Find([]byte(release))
		version := string(v)

		r := strings.Split(release, version)
		releaseWithoutVersion := r[0]
		releaseWithoutSuffix := strings.TrimSuffix(releaseWithoutVersion, "-")

		rvs[releaseWithoutSuffix] = string(version)
	}

	return rvs
}
