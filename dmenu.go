package dmenu

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

// Dmenu launches the dmenu command, piping in menuItems (separated by newlines)
// and passing provided args. The return value is the selection (which is whatever
// the user inputs in the case a menu item is not found), whether or not the
// selection is found in menuItems, and an error.
func Run(menuItems []string, args ...string) (selection string, found bool, err error) {
	cmd := exec.Command("dmenu", args...)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		return
	}

	// Pipe in the menu items
	for _, item := range menuItems {
		stdinPipe.Write([]byte(item + "\n"))
	}

	err = stdinPipe.Close()
	if err != nil {
		return
	}

	// Read the output of the command
	cmdOutput, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	// Trim the suffix
	selection = strings.TrimSuffix(string(cmdOutput), "\n")

	// Detect if the selection was in the menuItems
	for _, item := range menuItems {
		if item == selection {
			found = true
			break
		}
	}

	return
}
