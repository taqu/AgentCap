//go:build !unix

package agentbench

import "os/exec"

func isolateGroup(cmd *exec.Cmd) {}
