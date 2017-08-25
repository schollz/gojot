<p align="center">
<img
    src="https://raw.githubusercontent.com/schollz/gojot/v3/.github/gojot.png"
    width="260" height="80" border="0" alt="gojot">
<br>
<a href="https://github.com/schollz/gojot/releases/latest"><img src="https://img.shields.io/badge/version-3.0.0-brightgreen.svg?style=flat-square" alt="Version"></a>
<img src="https://img.shields.io/badge/coverage-40%25-yellow.svg?style=flat-square" alt="Code Coverage">
</p>

<p align="center">gojot is a modern command-line journal that is distributed and encrypted by default</p>

OK. But, really, *gojot* is just a fancy wrapper for `git`, `gpg` and `vim` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized in it a `git` repository.

Install
=======

First make sure you have `gpg`, `git`, and `vim` installed:

``` sourceCode
$ sudo apt-get install gpg git vim
```

Then you can install gojot using `pip`:

``` sourceCode
$ go get github.com/schollz/gojot
```

Usage
=====

First-time use
--------------

For the first time setup, just use:

    gojot

If you do not have any GPG keys you should first generate one with:

    gpg --gen-key