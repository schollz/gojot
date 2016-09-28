# sdees


[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
![](https://img.shields.io/badge/coverage-56.8%25-yellow.svg)

![Sdees](https://camo.githubusercontent.com/4e8c3264d9656edef93de9749f152b18314f780e/687474703a2f2f692e696d6775722e636f6d2f4936457a4544482e6a7067)

## SDEES is for Distributed Editing of Encrypted Stuff

Ok. But, really, `sdees` is just a fancy wrapper for `git` and `vim`/`nano`/`emacs` that allows you to make time-stamped entries to encrypted documents (like a notebook or journal) while keeping the entire document synchronized.

_Note_: The previous non-`git` version of `sdees` can [be found here](https://github.com/schollz/sdees/tree/1.X).


## Features

- Cross-compatibility (Windows/Linux/OS X).
- _Only two_ dependencies: `git` and a text editor (`vim` is bundled in Windows binary).
- Built-in encryption, compatible with `gpg`.
- Version control (all versions are saved, currently only newest is shown).
- Temp files are shredded (random bytes written before deletion).

## How it works

Each new entry in a document is symmetrically encrypted and inserted into a new orphan branch in the supplied `git` repository. The benefit of each entry having its own orphan branch is that each document will not have merge collisions when creating new entries on different local copies. Thus, `sdees` makes it perfectly safe to make new entries without internet access. Additional edits to an entry will commit a new (encrypted) change onto that document in the respective branch.

A single document is reconstructed by first fetching all remote branches, then filtering out which ones contain entries for the document of interest, decrypting each entry, and sorting the entries by date. Multiple documents can be stored in a single `git` repository. 

# Install

You can install by downloading [the latest release](https://github.com/schollz/sdees/releases/latest) or installing with Go 1.7+:
```
go get -u github.com/schollz/sdees
```

# Usage

The first time you run you can configure your remote system and editor.

```
sdees new.txt # edit a new document, new.txt
sdees --summary # list a summary
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
sdees --help # for more information
```

# Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT
