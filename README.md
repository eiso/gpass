# gpass

gpass is an encrypted account manager built on top of git. 

gpass creates a branch for each of the accounts you add, allowing you to have version control over your information such as passwords.

Commands under development: 

[] insert
[] show
[] list
[] rm
[] mv
[] cp
[] edit
[] generate
[] grep
[] help
[] version

gpass is inspired by [pass](https://www.passwordstore.org/), the Unix password manager, by [ZX2C4](https://www.zx2c4.com/). 

This project uses [go-git](https://github.com/src-d/go-git) for the git operations and [openpgp](https://godoc.org/golang.org/x/crypto/openpgp) for the encryption layer. 

This tool is currently under development. 
