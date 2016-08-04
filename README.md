![sdees](http://i.imgur.com/I6EzEDH.jpg)

A simple Go-program wrapper to **sync** a remote file, **decrypt** it, **edit** it, **encrypt** it, then **sync** it back.

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
```

One command instead of 6. One password instead of 3.

`sdees` also has cute features like:

* lite version control - it keeps track of diffs (encrypted)
* automatic date time (turn off with `--nodate`)
* list available files
* everything is always available locally (in `~/.sdees/`) accessed with `-l`

# Requirements

- `vim` or equivalent
- Have `someserver` that you already ran `ssh-copy-id someuser@someserver`

