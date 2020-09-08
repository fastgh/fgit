package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Debug ...
var Debug = false

// Mock ...
var Mock = false

// FileStat ...
func FileStat(path string, ensureExists bool) os.FileInfo {
	r, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(errors.Wrapf(err, "failed to stat file: %s", path))
		}
		if ensureExists {
			panic(errors.Wrapf(err, "file not exists: %s", path))
		}
		return nil
	}
	return r
}

// FileExists ...
func FileExists(path string) bool {
	if FileStat(path, false) == nil {
		return false
	}
	return true
}

// RemoveFile ...
func RemoveFile(path string) {
	if err := os.Remove(path); err != nil {
		panic(errors.Wrapf(err, "failed to delete file: %s", path))
	}
}

// ReadFile ...
func ReadFile(path string) []byte {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	return r
}

// WriteFile ...
func WriteFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0x777)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
}

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "failed to create directory: %s", path))
	}
}
