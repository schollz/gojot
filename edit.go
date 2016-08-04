package main

import (
	"fmt"
	"os"
	"os/exec"
)

// .vimrc
//
// func! WordProcessorModeCLI()
//     setlocal formatoptions=t1
//     setlocal textwidth=80
//     map j gj
//     map k gk
//     set formatprg=par
//     setlocal wrap
//     setlocal linebreak
//     setlocal noexpandtab
//     normal G$
// endfu
// com! WPCLI call WordProcessorModeCLI()

func editfile() {
	cmdArgs := []string{"-c", "WPCLI", "+startinsert", "new.txt"}
	cmd := exec.Command("vim", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(cmdArgs)
	fmt.Println(err)
}
