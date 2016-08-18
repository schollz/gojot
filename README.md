![build](https://img.shields.io/badge/build-passing-brightgreen.svg) ![Version 1.1.1](https://img.shields.io/badge/version-1.1.1-brightgreen.svg?version=flat-square)

![sdees](http://i.imgur.com/I6EzEDH.jpg)

# SDEES

...does **Serverless Decentralized Editing** of **Encrypted Stuff**.

...allows you to **Sync**, **Decrypt**, **Edit**, **Encrypt**, and **Sync** a document.

But, really, **SDEES** is just a fancy wrapper for `vim` that allows you to make time-stamped entries to an encrypted document (like a notebook or journal) while keeping the entire document synchronized remotely. Edits are stored individually and can be merged easily, so you can edit your document offline on multiple computers without worrying about losing data or overwriting.

This program grew out of constant utilization of `gpg`, `rsync`, and `vim`:

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

**SDEES** now combines this functionality into a single program (with `gpg` and `rsync` capabilities built-in):

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
- _Only one_ dependency: the text-editor `vim` (pre-bundled for Windows!).
- Encryption, compatible with `gpg`.
- Remote document transfer.
- Searching and summarizing,
- Version control (all versions are saved, but only newest is shown).

# Install

The simplest way to install is to just download the [latest release](https://github.com/schollz/sdees/releases/latest). To install from source you must install Go 1.6+.

```bash
git clone https://github.com/schollz/sdees.git
cd sdees
make install
```

Once installed you can update easily (must have Go1.6+ and Linux):

```bash
sdees --update
```

# Usage

The first time you run you can configure your remote system.

```bash
sdees new.txt # edit a new document, new.txt
sdees --summary -n 5 # list a summary of last five entries
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
```

![sdees usage](/branding/help2.gif)

# Acknowledgements

Logo design by [logodust](http://logodust.com)

Inspiration from [jrnl]

# License

MIT
