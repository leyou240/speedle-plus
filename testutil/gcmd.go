//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package testutil

import (
	"bytes"
	"os/exec"
)

// Cmd Class for exec.Cmd, standard output & error message & execute error
type Cmd struct {
	cmd      *exec.Cmd
	stdout   string
	stderr   string
	runError error
}

// Command Constructor for Cmd class
func Command(name string, arg ...string) *Cmd {
	return &Cmd{
		cmd: exec.Command(name, arg...),
	}
}

var cmdInstance = &Cmd{}

// Run execute the command.
func (c *Cmd) Run() {

	var outBuf bytes.Buffer
	c.cmd.Stdout = &outBuf

	var errBuf bytes.Buffer
	c.cmd.Stderr = &errBuf

	if err := c.cmd.Run(); err != nil {
		c.runError = err
	}

	c.stdout = string(outBuf.Bytes())
	c.stderr = string(errBuf.Bytes())
}

// Stdout Get the Stdout of exec.Cmd
func (c *Cmd) Stdout() string {

	return c.stdout
}

// Stderr Gget the Stderr of exec.Cmd
func (c *Cmd) Stderr() string {
	return c.stderr
}

// Success Return true if command run successfully
// Otherwise, return false
func (c *Cmd) Success() bool {
	return c.runError == nil
}
