package main

import (
	"os"
	"os/exec"
)

func main() {

	args := append([]string{"__complete"}, os.Args[1:]...)

	cmd := exec.Command("kubectl-konfigman", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}
