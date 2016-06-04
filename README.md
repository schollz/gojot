# droppybox

This is a Python3 wrapper for a simple file store.

# Requirements

- Python3.4+
- `gpg`, `rsync`
- Have `someserver` that you already ran `ssh-copy-id someuser@someserver`

# What does it do?

You can easily edit encrypted text files and store all the changes and keep the changes
in a server that will host copies for multiple computers to work on.

# Install

```bash
git clone https://github.com/schollz/droppybox.git && cd droppybox && sudo python3 setup.py install --record files.txt
```

# Uninstall

```bash
cat files.txt | xargs rm -rf
```

# Usage

```
$ z --help
usage: droppybox [-h] [-ls] [-l] [-e] [-u] [newfile]

positional arguments:
  newfile       work on a new file

optional arguments:
  -h, --help    show this help message and exit
  -ls, --list   list available files
  -l, --local   work locally
  -e, --edit    edit full document
  -u, --update  update droppybox
```
