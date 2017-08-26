// +build windows

package checkpath

import "path/filepath"

func pathsOnSameDevice(p1, p2 string) bool {
	id1 := filepath.VolumeName(p1)
	id2 := filepath.VolumeName(p2)
	return id1 == id2
}
