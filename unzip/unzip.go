package unzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Do(src, dest string) ([]string, error) {
	// Attmpeting to unzip downloaded tile into /tmp/grn/downloaded-tile/"
	fmt.Printf("Attempting to unzip: %s\ninto destination: %s\n", src, dest)

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// make folder to extract into
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// make file
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}

	}
	return filenames, nil
}
