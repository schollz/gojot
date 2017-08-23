# gogpg

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/schollz/gogpg)
![Coverage](https://img.shields.io/badge/coverage-82%25-green.svg?style=flat-square)

*gogpg* is a Go-library with a simple API for decrypting/encrypting single-user armor-encoded GPG files.

See the tests and documentation for more info.

## License

MIT

## To Do

- [ ] Add logging
- [ ] Automatically look for keyrings and use the first one found, otherwise throw error. Can override with specified key ring.
- [ ] Only specify the folder where the keyrings are
