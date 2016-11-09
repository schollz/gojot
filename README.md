## `sdees` is for `distributed` `editing` of `encrypted` `stuff`


[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
![](https://img.shields.io/badge/coverage-54%25-yellow.svg)


Ok. But, really, `sdees` is just a fancy wrapper for `git` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized. `sdees` is great for a *distributed and encrypted* journal and is compatible any local repo or hosted service (Gitlab/Bitbucket/Github). More info can be found in [INFO.md](https://github.com/schollz/sdees/blob/master/INFO.md).


_Note_: The previous non-`git` version of `sdees` can [be found here](https://github.com/schollz/sdees/tree/1.X).


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


## Install

First install the latest `git`. If you are on Windows/OSX just [use the build version](https://git-scm.com/downloads). If you are on Ubuntu, make sure to add the latest repo:

```
add-apt-repository ppa:git-core/ppa -y
apt-get update
apt-get install git -y
```

Then, to install `sdees`, simply download [the latest release binary](https://github.com/schollz/sdees/releases/latest)

_OR_

use `go get` if you have [installed Go](https://golang.org/dl/):

```
go get -u github.com/schollz/sdees
```

## Usage

```
sdees new.txt # edit a new document, new.txt
sdees --summary # list a summary
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
sdees --help # for more information
```

For more information, see https://sdees.schollz.com

### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT

This software includes third party open source software components: `src/base58*` and `src/merge*`. Each of these software components have their own license and copyright, stated in the respective files.
