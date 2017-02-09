# gojot is no longer supported, see the new, improved synchronized encrypted journal - [bol](https://github.com/schollz/bol)

<p align="center">
  <img src="https://gojot.schollz.com/_static/logo.png"/>
</p>

[![Version 2.1.0](https://img.shields.io/badge/version-2.1.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/gojot/releases/latest)
[![Github Releases](https://img.shields.io/github/downloads/schollz/gojot/latest/total.svg)](https://github.com/schollz/gojot/releases/latest)
[![Build Status](https://travis-ci.org/schollz/gojot.svg?branch=master)](https://travis-ci.org/schollz/gojot)
![](https://img.shields.io/badge/coverage-54%25-yellow.svg)
[![](https://img.shields.io/badge/gojot-documentation-blue.svg)](https://gojot.schollz.com/)

*gojot* is a modern command-line journal that is distributed and encrypted by default.

Ok. But, really, *gojot* is just a fancy wrapper for `git` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized. *gojot* is great for a *distributed and encrypted* journal and is compatible with all major operating systems. *gojot* is a single binary (with a text-editor built-in!), so you only need `git` to get started.

More info can be found in [the documentation](https://gojot.schollz.com/).

_Note_: The previous non-`git` version of *gojot* can [be found here](https://github.com/schollz/gojot/tree/1.X).


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

Example editing a [public Github repository](https://github.com/schollz/demo):

![](https://raw.githubusercontent.com/schollz/gojot/master/docs/source/_static/main_demo.gif)

```
gojot new.txt # edit a new document, new.txt
gojot --summary # list a summary
gojot --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
gojot --help # for more information
```
For more information, see https://gojot.schollz.com


## Requirements

First [install git](https://git-scm.com/downloads) (version 2.5+). If you are on Ubuntu, you can get the latest version of `git` using:
```
add-apt-repository ppa:git-core/ppa -y
apt-get update
apt-get install git -y
```

## Install

To install, just [download the latest *gojot* release](https://github.com/schollz/gojot/releases/latest).

_OR_

use `go get` if you have [installed Go](https://golang.org/dl/):

```
go get -u github.com/schollz/gojot
```


### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [these tools](https://gojot.schollz.com/about.html#alternatives-to-gojot).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT

This software includes third party open source software components: `src/base58*` and `src/merge*`. Each of these software components have their own license and copyright, stated in the respective files.
