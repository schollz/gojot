package main

const VIMRC = `set nocompatible
set backspace=2
func! WordProcessorModeCLI()
    setlocal formatoptions=t1
    setlocal textwidth=80
    map j gj
    map k gk
    set formatprg=par
    setlocal wrap
    setlocal linebreak
    setlocal noexpandtab
    normal G$
    normal zt
    set foldcolumn=7
    highlight Normal ctermfg=grey ctermbg=black
    hi NonText ctermfg=black guifg=black
endfu
com! WPCLI call WordProcessorModeCLI()
`

const VIMRC2 = `set nocompatible
set backspace=2
func! WordProcessorModeCLI()
    setlocal formatoptions=t1
    setlocal textwidth=80
    map j gj
    map k gk
    set formatprg=par
    setlocal wrap
    setlocal linebreak
    setlocal noexpandtab
    set foldcolumn=7
    normal G$
    highlight Normal ctermfg=grey ctermbg=black
    hi NonText ctermfg=black guifg=black
endfu
com! WPCLI call WordProcessorModeCLI()
`
