# go-pass

Is a native Go implementation of [Pass](https://www.passwordstore.org/), the unix password manager, by [ZX2C4](https://www.zx2c4.com/). It's currently under development. And looks to use [go-git](https://github.com/src-d/go-git) for the git operations and [openpgp](https://godoc.org/golang.org/x/crypto/openpgp) for the encryption layer. While pass is written in bash, using unix packages like `git` and `gpg2`, this project wants to increase the flexibility of how `pass` can be extended in Go. 
