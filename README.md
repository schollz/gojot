<p align="center">
  <img src="https://sdees.schollz.com/_static/logo.png"/>
</p>

[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
![](https://img.shields.io/badge/coverage-54%25-yellow.svg)
[![](https://img.shields.io/badge/sdees-documentation-blue.svg)](https://sdees.schollz.com/)

*sdees* is for *distributed editing* of *encrypted stuff*.

*sdees does editing, encryption,* and *synchronization*.

Ok. But, really, *sdees* is just a fancy wrapper for `git` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized. *sdees* is great for a *distributed and encrypted* journal and is compatible with all major operating systems. *sdees* is a single binary (with a text-editor built-in!), so you only need `git` to get started.

More info can be found in [the documentation](https://sdees.schollz.com/).

_Note_: The previous non-`git` version of *sdees* can [be found here](https://github.com/schollz/sdees/tree/1.X).


Features
--------
-  *Only one* dependency: `git` (version 2.5+)
-  Single binary, Cross-compatibility (Windows/Linux/OS X).
-  Fulltext encryption using OpenPGP, compatible with `gpg`.
-  Filename encryption using ChaCha20.
-  Built-in text editor, `micro` (but options for
   `emacs`/`vim`/`nano`).
-  Built-in version control (all versions are saved, currently only
   newest is shown).
-  Searching, summarizing, synchronized deletion, self-updating,
   collision management, and more.

## Usage

![](https://sdees.schollz.com/_static/main_demo.gif)
```
sdees new.txt # edit a new document, new.txt
sdees --summary # list a summary
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
sdees --help # for more information
```
For more information, see https://sdees.schollz.com


## Requirements

First [install git](https://git-scm.com/downloads) (version 2.5+). If you are on Ubuntu, you can get the latest version of `git` using:
```
add-apt-repository ppa:git-core/ppa -y
apt-get update
apt-get install git -y
```

## Install

To install, just [download the latest *sdees* release](https://github.com/schollz/sdees/releases/latest).

_OR_

use `go get` if you have [installed Go](https://golang.org/dl/):

```
go get -u github.com/schollz/sdees
```


### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT

This software includes third party open source software components: `src/base58*` and `src/merge*`. Each of these software components have their own license and copyright, stated in the respective files.
