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

	if len(os.Args) < 2 {
		fmt.Println("No args provided. See `grn -h`")
		os.Exit(1)
	}

	//rabbitmq-1.18/archive/p-rabbitmq-1.18.4-build.7.pivotal
	tileInS3 := flag.String("tile", "rabbitmq-1.18/archive/p-rabbitmq-1.18.4-build.7.pivotal", "tile path in s3 you want to generate release notes for.")

	flag.Parse()

	// create working directory
	path := "/tmp/grn/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	// cleanup
	defer rmDir(path)

	// Fetch tile
	tile, err := tilefetcher.Fetch(path, *tileInS3)
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

	outputHTML(releaseVersions)
}

func outputHTML(releaseVersions map[string]string) {
	for release, version := range releaseVersions {
		fmt.Println("<tr>")
		fmt.Printf("    <td>%s</td>\n", release)
		fmt.Printf("    <td>%s</td>\n", version)
		fmt.Println("</tr>")
	}
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
