## `sdees` is for `distributed` `editing` of `encrypted` `stuff`


[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
<!-- ![](https://img.shields.io/badge/coverage-45%25-yellow.svg) -->

Ok. But, really, `sdees` is just a fancy wrapper for `git` and `vim`/`nano`/`emacs` that allows you to make time-stamped entries to encrypted documents (like a notebook or journal) while keeping the entire document synchronized.

More info can be found in [INFO.md](https://github.com/schollz/sdees/blob/master/INFO.md).


_Note_: The previous non-`git` version of `sdees` can [be found here](https://github.com/schollz/sdees/tree/1.X).


## Features

- _Only two_ dependencies: `git` and a text editor (`vim` is bundled in Windows binary).
- Cross-compatibility (Windows/Linux/OS X).
- Built-in encryption, compatible with `gpg`.
- Version control (all versions are saved, currently only newest is shown).
- Temp files are shredded (random bytes written before deletion).
- Searching, summarizing, synchronized deletion, self-updating, and more.


## Install


Download [the latest release binary](https://github.com/schollz/sdees/releases/latest)

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

Here's `sdees` in action, editing a [demo repo](https://github.com/schollz/demo). Go check out the [demo repo](https://github.com/schollz/demo) and make sure its encrypted!

![](https://raw.githubusercontent.com/schollz/sdees/master/branding/help3.gif)


### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT
