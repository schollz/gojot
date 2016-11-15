<p align="center">
  <img src="https://jot.schollz.com/_static/logo.png"/>
</p>

[![Version 2.0.0beta](https://img.shields.io/badge/version-2.0.0beta-brightgreen.svg?version=flat-square)](https://github.com/schollz/jot/releases/latest)
[![Build Status](https://travis-ci.org/schollz/jot.svg?branch=master)](https://travis-ci.org/schollz/jot)
![](https://img.shields.io/badge/coverage-54%25-yellow.svg)
[![](https://img.shields.io/badge/jot-documentation-blue.svg)](https://jot.schollz.com/)

*jot* is a modern command-line journal that is distributed and encrypted by default.

Ok. But, really, *jot* is just a fancy wrapper for `git` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized. *jot* is great for a *distributed and encrypted* journal and is compatible with all major operating systems. *jot* is a single binary (with a text-editor built-in!), so you only need `git` to get started.

More info can be found in [the documentation](https://jot.schollz.com/).

_Note_: The previous non-`git` version of *jot* can [be found here](https://github.com/schollz/jot/tree/1.X).


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
![](https://raw.githubusercontent.com/schollz/jot/master/docs/source/_static/main_demo.gif)

```
jot new.txt # edit a new document, new.txt
jot --summary # list a summary
jot --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
jot --help # for more information
```
For more information, see https://jot.schollz.com


## Requirements

First [install git](https://git-scm.com/downloads) (version 2.5+). If you are on Ubuntu, you can get the latest version of `git` using:
```
add-apt-repository ppa:git-core/ppa -y
apt-get update
apt-get install git -y
```

## Install

To install, just [download the latest *jot* release](https://github.com/schollz/jot/releases/latest).

_OR_

use `go get` if you have [installed Go](https://golang.org/dl/):

```
go get -u github.com/schollz/jot
```


### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [these tools](https://jot.schollz.com/about.html#alternatives-to-jot).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT

This software includes third party open source software components: `src/base58*` and `src/merge*`. Each of these software components have their own license and copyright, stated in the respective files.
