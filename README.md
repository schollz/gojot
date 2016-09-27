# gitsdees


[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-brightgreen.svg?version=flat-square)](https://github.com/schollz/gitsdees/releases/latest)
[![Build Status](https://travis-ci.org/schollz/gitsdees.svg?branch=master)](https://travis-ci.org/schollz/gitsdees)
![](https://img.shields.io/badge/coverage-56.8%25-yellow.svg)

## SDEES is for Distributed Editing of Encrypted Stuff

Ok. But, really, `sdees` is just a fancy wrapper for `git` and `vim`/`nano`/`emacs` that allows you to make entries to an encrypted document (like a notebook or journal) while keeping the entire document synchronized.

This program grew out of constant utilization of `gpg`, `rsync`, and `vim`/`nano`/`emacs`. Before `sdees` I had to do this:

```
$ rsync -arq --update user@remote:encrypted_notes encrypted_notes
$ gpg -d encrypted_notes > notes
Enter passphrase: *******
$ vim notes
$ gpg --symmetric notes -o encrypted_notes
Enter passphrase: *******
Repeat passphrase: *******
File encrypted_notes exists. Overwrite? (y/N) y
$ rm notes
$ rsync -arq --update encrypted_notes user@remote:encrypted_notes
```

Now, with `sdees` (which has `gpg` and `rsync` capabilities built-in) I can do this:

```
$ sdees
Pulling from remote...done.
Enter password for editing 'notes.txt': ******
Wrote M6nWWLw.fNfEqs.0qPJJeZ.gpg.
+410 words.
Pushing to remote...done.
```

More information [in the code](https://github.com/schollz/sdees/blob/master/main.go#L1-L29).

## Features

- Cross-compatibility (Windows/Linux/OS X).
- _Only two_ dependencies: `git` and a text editor (`vim` is bundled in Windows binary).
- Built-in encryption, compatible with `gpg`.
- Version control (all versions are saved, currently only newest is shown).
- Temp files are shredded (random bytes written before deletion).

# Install

The simplest way to install is to just download the [latest release](https://github.com/schollz/sdees/releases/latest). To install from source you must install Go 1.6+.

```
go get -u github.com/schollz/gitsdees
```

# Usage

The first time you run you can configure your remote system and editor.

```
sdees new.txt # edit a new document, new.txt
sdees --summary -n 5 # list a summary of last five entries
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
sdees --help # for more information
```

![sdees usage](/branding/help2.gif)

# Acknowledgements

Logo graphic from [logodust](http://logodust.com).

Inspiration from [jrnl](http://jrnl.sh/).

Southwest Airlines for providing two mechanical failures that gave me 8+ extra hours to code this.

Stack overflow (see code for attributions).

# License

MIT
