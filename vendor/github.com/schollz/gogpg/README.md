# gogpg

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/schollz/gogpg)
![Coverage](https://img.shields.io/badge/coverage-76%25-green.svg?style=flat-square)

*gogpg* is a Go-library with a simple API for decrypting/encrypting single-user armor-encoded GPG files.

See the tests and documentation for more info.

## Development

To run the tests, you need to generate a key for `Testy McTestFace`:

```
$ cd testing
$ gpg --gen-key
$ # Use ID "Testy McTestFace" and password "1234"
$ gpg --yes --armor --recipient "Testy McTestFace" --trust-model always --encrypt hello.txt
```

## License

MIT
