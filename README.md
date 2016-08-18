![build](https://img.shields.io/badge/build-passing-brightgreen.svg) ![Version 1.1.1](https://img.shields.io/badge/version-1.1.1-brightgreen.svg?version=flat-square)

![sdees](http://i.imgur.com/I6EzEDH.jpg)

# SDEES

- ...is a program that allows **Serverless Decentralized Editing of Encrypted Stuff**.
- ...is for **Syncing** remote documents, **Decrypting**, **Editing**, **Encrypting**, then **Syncing** back.

But really, **SDEES** is just a fancy wrapper for `vim` that allows you to make time-stamped entries to an encrypted document (like a notebook or journal) while keeping the entire document in sync remotely. Since all changes are stored individually they can easily be merged, which means that you can edit your document offline on multiple computers without worrying about merging or overwriting.

This program grew out of constant utilization of `gpg`, `rsync`, and `vim`:

```bash
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

which **SDEES** now combines into the single command (with `gpg` and `rsync` capabilities baked in now):

```bash
$ sdees
Pulling from remote...done.
Enter password for editing 'notes.txt': ******
Wrote M6nWWLw.fNfEqs.0qPJJeZ.gpg.
+410 words.
Pushing to remote...done.
```

More information [in the code](https://github.com/schollz/sdees/blob/master/main.go#L1-L39).

## Features

- Cross-compatibility (Windows/Linux/OS X)
- _Only one_ dependency: the text-editor `vim` (pre-bundled for Windows!)
- Encryption, compatible with `gpg`
- Remote document transfer
- Searching and summarizing
- Version control (all versions are saved, but only newest is shown)

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

# Acknowledgements

Logo design by [logodust](http://logodust.com)

# License

MIT
