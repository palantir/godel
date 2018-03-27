package installupdate

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// RunUpgradeLegacyConfig runs the "upgrade-config" task in legacy mode by invoking
// "{{projectDir}}/godelw upgrade-config --legacy". Sets the sets the Stdout and Stderr to that of the current process.
func RunUpgradeLegacyConfig(projectDir string) error {
	godelw := path.Join(projectDir, "godelw")
	cmd := exec.Command(godelw, "upgrade-config", "--legacy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// MajorVersion returns the major version returned by "{{projectDir}}/godelw version" (if, for "x.y.z+", "x" parses as
// an integer).
func MajorVersion(projectDir string) (int, error) {
	version, err := projectVersionUsingGodelw(projectDir)
	if err != nil {
		return -1, err
	}
	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		return -1, errors.Wrapf(err, "version does not consist of at least 3 '.'-delimited parts")
	}
	versionInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return -1, errors.Wrapf(err, "failed to parse major version as integer")
	}
	return versionInt, nil
}

func projectVersionUsingGodelw(projectDir string) (string, error) {
	godelw := path.Join(projectDir, "godelw")
	cmd := exec.Command(godelw, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute command %v: %s", cmd.Args, string(output))
	}
	outputString := strings.TrimSpace(string(output))
	parts := strings.Split(outputString, " ")
	if len(parts) != 3 {
		return "", errors.Errorf(`expected output %s to have 3 parts when split by " ", but was %v`, outputString, parts)
	}
	return parts[2], nil
}
