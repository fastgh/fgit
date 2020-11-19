package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Debug ...
var Debug = false

// FileStat ...
func FileStat(path string, ensureExists bool) os.FileInfo {
	r, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(errors.Wrapf(err, "获取文件信息失败: %s", path))
		}
		if ensureExists {
			panic(errors.Wrapf(err, "文件不存在: %s", path))
		}
		return nil
	}
	return r
}

// FileExists ...
func FileExists(path string) bool {
	fi := FileStat(path, false)
	if fi == nil {
		return false
	}
	if fi.IsDir() {
		panic(fmt.Errorf("%s是目录，不是文件", path))
	}
	return true
}

// DirExists ...
func DirExists(path string) bool {
	fi := FileStat(path, false)
	if fi == nil {
		return false
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("%s是文件，不是目录", path))
	}
	return true
}

// RemoveFile ...
func RemoveFile(path string) {
	if err := os.Remove(path); err != nil {
		panic(errors.Wrapf(err, "删除文件失败: %s", path))
	}
}

// ReadFile ...
func ReadFile(path string) []byte {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrapf(err, "读取文件失败: %s", path))
	}
	return r
}

// WriteFile ...
func WriteFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0x777)
	if err != nil {
		panic(errors.Wrapf(err, "写入文件失败: %s", path))
	}
}

// MkdirAll ...
func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(errors.Wrapf(err, "创建目录失败: %s", path))
	}
}

// JSONPretty ...
func JSONPretty(data interface{}) string {
	json, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(errors.Wrapf(err, "JSON序列化失败"))
	}
	return string(json)
}

// JSONMarshal ...
func JSONMarshal(data interface{}) string {
	json, err := json.Marshal(data)
	if err != nil {
		panic(errors.Wrapf(err, "JSON序列化失败"))
	}
	return string(json)
}

// JSONUnmarshal ...
func JSONUnmarshal(jsonText string, v interface{}) {
	if err := json.Unmarshal([]byte(jsonText), v); err != nil {
		panic(errors.Wrapf(err, "JSON反序列化失败: %s\n", jsonText))
	}
}

// GetWorkDir ...
func GetWorkDir() string {
	r, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "无法获取当前工作目录路径"))
	}
	return r
}
