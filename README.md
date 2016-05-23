FSRename v2 - Fast & Simple Rename
================================

[![Build Status](https://travis-ci.org/c9s/fsrename.svg?branch=master)](https://travis-ci.org/c9s/fsrename)

A simple, powerful rename tool supports complex filtering

fsrename separates the pattern option, therefore you can specify the pattern without typeing these character escapes.

further, this tool supports concurrent renaming (powered by Go's channel and routine)

Install
--------------

    go get -u -x github.com/c9s/fsrename/fsrename

Usage
---------------------

    fsrename [options] [path...]

> Note: When `[path...]` is not given, "./" will be used as the default path for scanning files.

To see the documentation in console:

    go doc github.com/c9s/fsrename/fsrename

You can link the binary file to your GOPATH to create an doc alias for this

    P=$(echo $GOPATH|cut -d: -f1)
    ln -s $P/src/github.com/c9s/fsrename/fsrename $P/src/fsrename

## Options

### Filter Options

`-match` pre-filter the files and directories based on the given regular pattern.

`-contains` pre-filter the files and directories based on the given string needle.

`-file` only for files.

`-f` an alias of `-file`

`-dir` only for directories.

`-d` an alias of `-d`

`-ext` find files with matched file extension.


### Replacement Options

Please note the replacement target only works for the basename of a path.
`-replace*` and `-with*` should be combined together to replace the substrings.

`-replace` specify target substring with normal string matching.
`-r` alias of `-replace`


`-replaceRegexp` specify target substring with regular expression matching.
`-rre` alias of `-replaceRegexp`


`-with` replacement for the target substring.
`-w` alias of `-with`


`-withFormat` replacement with fmt.Sprintf format for the target substring.

### Replace Rule Builder Options

`-trimPrefix` trim filename prefix.

`-trimSuffix` trim filename suffix (this option removes suffix even for filename extensions).

`-camel` converts dash/underscore separated filenames into camelcase filenames.

`-underscore` converts camelcase filesnames into underscore separated filenames.

### Common Options

`-dryrun`  dry run, don't rename, just preview the result.




Quick Examples
-------------

Find files with extension `.php` and replace the substring from the filename.

    fsrename -ext "php" -replace "some" -with "others" src/

Replace `Stmt.go` with "_stmt.go" under the current directory:

    fsrename -replace "Stmt.go" -with "_stmt.go"

Replace `Stmt.go` with "_stmt.go" under directory `src/c6`:

    fsrename -replace "Stmt.go" -with "_stmt.go" src/c6

Replace `foo` with `bar` from files contains `prefix_` 

    fsrename -file -contains prefix_ -replace foo -with bar test

Or use `-match` to pre-filter the files with regular expression

    fsrename -file -match "[a-z]+" -replace foo -with bar test

Use regular expression without escaping:

    fsrename -replaceRegexp "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -replaceRegexp "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -ext go -replaceRegexp "[a-z]*" -with "123" src/c6

    fsrename -dir -replace "_xxx" -with "_aaa" src/c6

    fsrename -replace "_xxx" -with "_aaa" -dryrun  src/c6



## Roadmap

v2.1

- [ ] Add `-list` to print the filtered file paths instead of renaming the files.
- [ ] Support more actions rather than file rename.
  - [ ] `-backup` to backup selected files.
- [ ] Add rename log printer to support rollback.
- [ ] Add `-grep` to grep files that contains specific pattern.


## LICENSE

MIT License

