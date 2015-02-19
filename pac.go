package pac

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/getlantern/byteexec"
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("pac")

	mu sync.Mutex
	be *byteexec.Exec
)

// EnsureHelperToolPresent checks if helper tool exists and extracts it if not.
// On Mac OS, it also checks and set the file's owner to root:wheel and the setuid bit,
// it will request user to input password through a dialog to gain the rights to do so.
// fullPath: the file to be checked and generated if not exists.
// prompt: the message to be shown on the dialog.
// iconPath: the full path of the icon to be shown on the dialog.
func EnsureHelperToolPresent(fullPath string, prompt string, iconFullPath string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	pacBytes, err := Asset("pac")
	if err != nil {
		return fmt.Errorf("Unable to access pac asset: %v", err)
	}
	be, err = byteexec.New(pacBytes, fullPath)
	if err != nil {
		return fmt.Errorf("Unable to extract helper tool: %v", err)
	}
	return ensureElevatedOnDarwin(be, fullPath, prompt, iconFullPath)
}

/* On tells OS to configure proxy through `pacUrl` */
func On(pacUrl string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}

	log.Debug(be.Filename)
	cmd := be.Command("on", pacUrl)
	return run(cmd)
}

/* Off sets proxy mode back to direct/none */
func Off() (err error) {
	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}
	cmd := be.Command("off")
	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Unable to execute helper tool: %s\n%s", err, string(out))
	}
	log.Tracef("Command %v output %v", cmd.Path, string(out))
	return nil
}
