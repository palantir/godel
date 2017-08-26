// +build linux darwin

package checkpath

import (
	"os"
	"syscall"
)

func pathsOnSameDevice(p1, p2 string) bool {
	id1, ok := getDeviceID(p1)
	if !ok {
		return false
	}
	id2, ok := getDeviceID(p2)
	if !ok {
		return false
	}
	return id1 == id2
}

func getDeviceID(p string) (interface{}, bool) {
	fi, err := os.Stat(p)
	if err != nil {
		return 0, false
	}
	s := fi.Sys()
	switch s := s.(type) {
	case *syscall.Stat_t:
		return s.Dev, true
	}
	return 0, false
}
