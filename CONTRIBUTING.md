### How to contribute

1. Fork [https://github.com/DE-labtory/koa](https://github.com/DE-labtory/koa)
2. Register issue you are going to work, and discuss with maintainers whether this feature is needed. Or you can choose existing issue.
3. After assigned the issue, work on it
4. Document about your code
5. Test with `go test -v ./...`
6. Format your code with `goimports -w ./` commands
7. After passing tc and formatting, create pr to `develop` branch referencing issue number
8. Passing travis-ci with more than one approve from maintainers, you can **rebase and merge** to the develop branch
9. After passing all the tc, no build error we can merge to `master` branch at milestone

### Rules to manage branch

* `master`: project with release level can be merged to `master` branch
* `develop`: new feature developed after fully verified from others can be merged to `develop` branch

### Tip

* For not overlapping work with others, please work on small feature as possible.
* When register new issue, concrete documenting issue helps other to understand what you are going to work on and to feedback about your proposal easily.
* Start your commit message with capital letter.  