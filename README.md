![build](https://img.shields.io/badge/build-passing-brightgreen.svg)

![sdees](http://i.imgur.com/I6EzEDH.jpg)

**SDEES** is for **syncing** remote files, **decrypting**, **editing**, **encrypting**, then **syncing** back.

However, **SDEES** is also a program that allows **serverless decentralized editing of encrypted stuff**.

That is, you can use it offline/online and never fear of losing data or having trouble merging encrypted edits.

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

**SDEES** lets you do this:

```bash
$ sdees
Enter password for editing: ******
```

One command instead of 6\. One password instead of 3.

# Requirements

- `vim` or equivalent
- Have `someserver` that you already ran `ssh-copy-id someuser@someserver`

# Install

To install, you must install Go 1.6+.

```
make
```

# Run

To run, just use

```
./sdees
```

For more information use

```
./sdees --help
```
