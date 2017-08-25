<p align="center">
<img
    src="https://raw.githubusercontent.com/schollz/gojot/v3/.github/gojot.png"
    width="260" height="80" border="0" alt="gojot">
<br>
<a href="https://github.com/schollz/gojot/releases/latest"><img src="https://img.shields.io/badge/version-3.0.0-brightgreen.svg?style=flat-square" alt="Version"></a>
<img src="https://img.shields.io/badge/coverage-40%25-yellow.svg?style=flat-square" alt="Code Coverage">
</p>

<p align="center">gojot is a modern command-line journal that is distributed and encrypted by default</p>

OK. But, really, *gojot* is just a fancy wrapper for `git`, `gpg` and `vim` that allows you to make time-stamped entries to encrypted documents while keeping the entire document synchronized in it a `git` repository.

Install
=======

First make sure you have `gpg`, `git`, and `vim` installed:

``` sourceCode
$ sudo apt-get install gpg git vim
```

Then you can install gojot using `pip`:

``` sourceCode
$ go get github.com/schollz/gojot
```

Usage
=====

First-time use
--------------

For the first time setup, just use:

    gojot

If you do not have any GPG keys you should first generate one with:

    gpg --gen-key

License
========

MIT

History
=======

Version 4 (current version) ([5464ef97](https://github.com/schollz/gojot/tree/5464ef97c3983b994072d3737f235f35b698b48e))

- Finished August 25th, 2017
- 1,172 lines of Go
- Requires `gpg` for encryption
- Requires `git` for syncing


Version 3 ([5faaeb3](https://github.com/schollz/gojot/tree/5faaeb3))

- Finished August 8th, 2017
- 495 lines of Python
- Requires `gpg` for encryption
- Requires `git` for syncing

Version 2 ([f881b416](https://github.com/schollz/gojot/tree/f881b416))

- Finished November 17th, 2016
- 4,163 lines of Go
- Built-in `gpg` for encryption
- Requires `git` for syncing


Version 1 ([03b4419a](https://github.com/schollz/gojot/tree/03b4419a1e9a032db8dd96cd18517a0830db4626))

- Finished October 3rd, 2016
- 1,633 lines of Go
- Built-in `gpg` for encryption
- Built-in `rsync` for syncing 

Version 0 ([d6b66c3c](https://github.com/schollz/gojot/tree/d6b66c3c1ac7fa6e34b971342c0f5257e8b7af30))

- Finished August 3rd, 2016
- 341 lines of Python.
- Requires `gpg` for encryption
- Requires `rsync` for syncing

