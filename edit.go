package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

// .vimrc
//

func editfile() {
	logger.Debug("Editing file")

	err := ioutil.WriteFile(path.Join(RuntimeArgs.TempPath, "vimrc"), []byte(`func! WordProcessorModeCLI()
    setlocal formatoptions=t1
    setlocal textwidth=80
    map j gj
    map k gk
    set formatprg=par
    setlocal wrap
    setlocal linebreak
    setlocal noexpandtab
    normal G$
endfu
com! WPCLI call WordProcessorModeCLI()`), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmdArgs := []string{"-u", path.Join(RuntimeArgs.TempPath, "vimrc"), "-c", "WPCLI", "+startinsert", path.Join(RuntimeArgs.TempPath, "temp")}
	cmd := exec.Command("vim", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	encrypt()
}
