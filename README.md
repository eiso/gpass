# gpass

gpass is an encrypted account manager built on top of git. 

gpass creates a branch for each of the accounts you add, allowing you to have version control over your encrypted information such as passwords.

Commands under development: 

- [x] help
- [x] version
- [x] init
  - [x] existing repository
  - [x] existing private key
  - [ ] create new repository
  - [ ] create new private key
- [x] insert
  - [x] single line
  - [ ] multiple line (editor)
- [x] show
- [x] list
- [x] rm
- [x] mv
- [x] cp
- [ ] edit
- [ ] generate
- [ ] search
- [ ] grep


gpass is inspired by [pass](https://www.passwordstore.org/), the Unix password manager, by [ZX2C4](https://www.zx2c4.com/). 

This project uses [go-git](https://github.com/src-d/go-git) for the git operations and [openpgp](https://godoc.org/golang.org/x/crypto/openpgp) for encryption. 

This tool is currently under development. 
