package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

// This is an example of a script you could
// write to periodically compact your database.

// IMPORTANT!!!  Before running this script or
// one like it that you write, make sure no
// processes are using the database!!!

const dirPath = "../example_data/"
const tblExt = ".json"

func main() {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := file.Name()

		if filename == "mst3k_episodes_default.json" {
			continue
		}

		// If entry is sub dir, current dir, or parent dir, skip it.
		if file.IsDir() || filename == "." || filename == ".." {
			continue
		}

		if !strings.HasSuffix(filename, tblExt) {
			continue
		}

		compactFile(filename)
		if err != nil {
			panic(err)
		}
	}
}

func compactFile(filename string) error {
	filepath := dirPath + filename
	backupFilepath := dirPath + strings.TrimSuffix(filename, tblExt) + ".old"

	// Move the table to a backup file.
	if err := os.Rename(filepath, backupFilepath); err != nil {
		return err
	}

	oldfile, err := os.Open(backupFilepath)
	if err != nil {
		return err
	}
	defer oldfile.Close()

	newfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer newfile.Close()

	oldfileScanner := bufio.NewScanner(oldfile)
	for oldfileScanner.Scan() {
		str := strings.TrimRight(oldfileScanner.Text(), "X")

		if len(str) > 0 {
			_, err := newfile.WriteString(str + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err := oldfileScanner.Err(); err != nil {
		return err
	}

	newfile.Sync()

	return nil
}
