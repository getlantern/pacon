// +build !darwin

package pac

import (
	"github.com/getlantern/byteexec"
)

func ensureElevatedOnDarwin(be *byteexec.Exec, helperFullPath string, prompt string, iconFullPath string) (err error) {
	return nil
}
