# sdees

A simple Python3 wrapper to **sync** a remote file, **decrypt** it, **edit** with vim, **encrypt** it, then **sync** it back.

# About

Instead of doing this:

```
$ rsync -arq --update user@remote:encryptedfile encryptedfile
$ gpg -d encryptedfile > file
Enter passphrase: *******
$ vim file
$ gpg --symmetric file -o encryptedfile
Enter passphrase: *******
Repeat passphrase: *******
File `encryptedfile' exists. Overwrite? (y/N) y
$ rm file
$ rsync -arq --update encryptedfile user@remote:encryptedfile
```

`sdees` lets you do this:

```bash
$ sdees
Enter passphrase: *******
```

One command instead of 6. One password instead of 3. `sdees` also has lite version control - it keeps track of diffs (encrypted) so that every change could be recapitulated.

# Requirements

- Python3.4+
- `gpg`, `rsync`
- Have `someserver` that you already ran `ssh-copy-id someuser@someserver`

# What does it do?

You can easily edit encrypted text files and store all the changes and keep the changes
in a server that will host copies for multiple computers to work on.

# Install

```bash
git clone https://github.com/schollz/sdees.git && cd sdees && sudo python3 setup.py install --record files.txt
```

# Uninstall

```bash
cat files.txt | xargs rm -rf
```

# Usage

```
$ sdees --help
usage: sdees [-h] [-ls] [-l] [-e] [-u] [newfile]

positional arguments:
  newfile       work on a new file

optional arguments:
  -h, --help    show this help message and exit
  -ls, --list   list available files
  -l, --local   work locally
  -e, --edit    edit full document
  -u, --update  update droppybox
```
