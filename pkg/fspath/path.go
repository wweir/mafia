package fspath

import (
	"bytes"
	"strings"
)

type pathRelation int8

const (
	PathIrrelevant pathRelation = iota
	PathSup
	PathParrent
	PathSelf
	PathChild
	PathSub
)

func SumPathRelation(self, target string) pathRelation {
	self = strings.TrimPrefix(self, ".")
	target = strings.TrimPrefix(target, ".")
	selfPath := make([]byte, 0, len(self)+2)
	targetPath := make([]byte, 0, len(target)+2)

	if self == "" || self[0] != '/' {
		selfPath = append(selfPath, '/')
	}
	selfPath = append(selfPath, []byte(self)...)
	if selfPath[len(selfPath)-1] != '/' {
		selfPath = append(selfPath, '/')
	}

	if target == "" || target[0] != '/' {
		targetPath = append(targetPath, '/')
	}
	targetPath = append(targetPath, []byte(target)...)
	if targetPath[len(targetPath)-1] != '/' {
		targetPath = append(targetPath, '/')
	}

	if bytes.Equal(selfPath, targetPath) {
		return PathSelf
	}

	if bytes.HasPrefix(selfPath, targetPath) {
		left := selfPath[len(targetPath):]
		if bytes.IndexByte(left[:len(left)-1], '/') < 0 {
			return PathParrent
		}
		return PathSup
	}

	if bytes.HasPrefix(targetPath, selfPath) {
		left := targetPath[len(selfPath):]
		if bytes.IndexByte(left[:len(left)-1], '/') < 0 {
			return PathChild
		}
		return PathSub
	}

	return PathIrrelevant
}

func SplitPrefixDir(path string) (prefix, last string) {
	switch path {
	case "", ".", "/", "./":
		return "", ""
	}

	if path[0] == '/' {
		path = path[1:]
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if idx := strings.IndexByte(path, '/'); idx > 0 {
		return path[:idx], path[idx+1:]
	}
	return path, ""
}

func Relative(path string) string {
	if path == "" || path == "/" {
		return "."
	}
	if path[0] == '/' {
		return path[1:]
	}
	return path
}
