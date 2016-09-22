[![Build Status](https://travis-ci.org/schollz/sdees.svg?branch=master)](https://travis-ci.org/schollz/sdees)
[![Version 1.2.5](https://img.shields.io/badge/version-1.2.5-brightgreen.svg?version=flat-square)](https://github.com/schollz/sdees/releases/latest)

![sdees](http://i.imgur.com/I6EzEDH.jpg)

## SDEES Does Editing, Encryption, and Synchronization

Ok. But, really, `sdees` is just a fancy wrapper for `vim`/`nano`/`emacs` that allows you to make time-stamped entries to an encrypted document (like a notebook or journal) while keeping the entire document synchronized remotely.

The remote computer is used only for file storage and does not require any server-side code. Edits are stored individually and can be merged easily, so you can edit your document offline on multiple computers without worrying about losing data or overwriting.

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
- _Only one_ dependency: a CLI text-editor, either `vim`, `nano`, or `emacs`. (The Windows release always includes `vim.exe`).
- Encryption, compatible with `gpg`.
- Remote document transfer.
- Searching and summaries.
- Version control (all versions are saved, currently only newest is shown).
- Temp files are shredded (random bytes written before deletion).

# Install

The simplest way to install is to just download the [latest release](https://github.com/schollz/sdees/releases/latest). To install from source you must install Go 1.6+.

```
git clone https://github.com/schollz/sdees.git
cd sdees
make
make install    # if you want to install to your path
```

Once installed you can update easily (must have Go1.6+ and Linux):

```
sdees --update
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
