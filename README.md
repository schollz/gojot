![build](https://img.shields.io/badge/build-passing-brightgreen.svg)

![sdees](http://i.imgur.com/I6EzEDH.jpg)

# SDEES

- ...is a program that allows **Serverless Decentralized Editing of Encrypted Stuff**.
- ...is for **Syncing** remote files, **Decrypting**, **Editing**, **Encrypting**, then **Syncing** back.

Who am I kidding? **SDEES** is a really just a fancy wrapper for `vim` that allows you to make time-stamped entries to an encrypted document (like a notebook or journal) while keeping the entire document in sync remotely. Since all changes are stored individually they can easily be merged. That means that you can edit your document offline on multiple computers without worrying about merging those changes later.

# About

Instead of doing this:

```
$ rsync -arq --update user@remote:encrypted_notes encrypted_notes
$ gpg -d encrypted_notes > notes
Enter passphrase: *******
$ vim notes
$ gpg --symmetric notes -o encrypted_notes
Enter passphrase: *******
Repeat passphrase: *******
File `encrypted_notes' exists. Overwrite? (y/N) y
$ rm notes
$ rsync -arq --update encrypted_notes user@remote:encrypted_notes
```

**SDEES** lets you do this:

```bash
$ sdees
Enter password for editing: ******
```

One command instead of 6\. One password instead of 3\. And no worries about overwriting and losing data.

## Features

- _Only one_ dependency: the text-editor `vim` (pre-bundled for Windows!)
- GPG-based encryption
- Builtin remote file transfer
- Searching and summarizing
- Cross-compatibility (Windows/Linux/OS X)

# Install

The simplest way to install is to just download the [latest release](https://github.com/schollz/sdees/releases/latest). To install from source you must install Go 1.6+.

```
git clone https://github.com/schollz/sdees.git
cd sdees
make install
```

# Usage

The first time you run you can configure your remote system.

```bash
sdees new.txt # edit a new document, new.txt
sdees --summary -n 5 # list a summary of last five entries
sdees --search "dogs cats" # find all entries that mention 'dogs' or 'cats'`
```
