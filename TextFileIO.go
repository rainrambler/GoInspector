package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ReadLines reads a whole file into memory and returns a slice of its lines
// for whole file as a string, use ioutil.ReadFile instead.
func ReadLines(fullpath string) ([]string, error) {
	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//WriteLines - each line should NOT end with \n or \r\n
func WriteLines(lines []string, fullpath string) error {
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}

func IoReader(fullpath string) io.ReaderAt {
	r, err := os.Open(fullpath)
	if err != nil {
		panic(err)
	}
	return r
}

func ReadBinFile(fullpath string) ([]byte, error) {
	return ioutil.ReadFile(fullpath)
}
