# sdees


[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
![](https://img.shields.io/badge/coverage-53%25-yellow.svg)


## SDEES is for Distributed Editing of Encrypted Stuff

Ok. But, really, `sdees` is just a fancy wrapper for `git` and `vim`/`nano`/`emacs` that allows you to make time-stamped entries to encrypted documents (like a notebook or journal) while keeping the entire document synchronized.

_Note_: The previous non-`git` version of `sdees` can [be found here](https://github.com/schollz/sdees/tree/1.X).


## Features

- _Only two_ dependencies: `git` and a text editor (`vim` is bundled in Windows binary).
- Cross-compatibility (Windows/Linux/OS X).
- Built-in encryption, compatible with `gpg`.
- Version control (all versions are saved, currently only newest is shown).
- Temp files are shredded (random bytes written before deletion).


## Install

You can install by downloading [the latest release](https://github.com/schollz/sdees/releases/latest) or installing with Go 1.7+:
```
go get -u github.com/schollz/sdees/...
```

Info about using a `git` repository other than Github/Bitbucket can [be found in INFO.md](https://github.com/schollz/sdees/blob/master/INFO.md).


## How it works

When `sdees` starts for the first time it will request a `git` repository which can be [local or remote](https://github.com/schollz/sdees/blob/master/INFO.md). Then `sdees` provides an option to write text using the editor of choice into a new *entry*. Each new entry is inserted into a new orphan branch in the supplied `git` repository. The benefit of placing each entry into its own orphan branch is that merge collisions are avoided after creating new entries on different machines. Thus, `sdees` makes it perfectly safe to make *new* entries without internet access.

The combination of entries be displayed as a *document*. The document is reconstructed by 1) fetching all remote branches, 2) filtering out which ones contain entries for the document of interest, 3) decrypting each entry (if needed), and then 4) sorting the entries by the commit date. Merge conflicts should only occur when simultaneous edits are made to the same entry (which is not the use case here, in general). Multiple documents can be stored in a single `git` repository.

Optionally, all information saved in the `git` repo can be encrypted using a symmetric cipher with a user-provided passphrase. The passphrase is not stored anywhere on the machine or repo. When enabled, each entry in the `git` repo is encrypted. When editing an encrypted document, a decrypted temp file is stored and then shredded (random bytes written and then deleted) after use.

## Usage

The first time you run you can configure your remote system and editor.

```
sdees new.txt # edit a new document, new.txt
sdees --summary # list a summary
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
sdees --help # for more information
```


### Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT
