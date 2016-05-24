FSRename v2 - Fast & Simple Rename
================================

[![Build Status](https://travis-ci.org/c9s/fsrename.svg?branch=master)](https://travis-ci.org/c9s/fsrename)

A simple, powerful rename tool supports complex filtering, written in GO.

`fsrename` separates the pattern/replace options, therefore you can specify the
pattern without typing these escaping characters and the slashes.

pre-filtering, extension filtering, prefix trimming, suffix trimming ,
camelcase conversion, underscore conversion are all supported.

> This tool is different from `gorename`, `gorename` is for refactoring your code, rename variables.


SYNOPSIS
--------------

    fsrename -file -match "[a-z]+" -rr "foo_+" -with bar test

INSTALL
--------------

    go get -u -x github.com/c9s/fsrename/fsrename

USAGE
---------------------

    fsrename [options] [path...]

> Note: When `[path...]` is not given, "./" will be used as the default path for scanning files.

To see the documentation in console:

    go doc github.com/c9s/fsrename/fsrename

You can create a link to the package under your GOPATH to create an doc alias

    ln -s $(go list -f "{{.Dir}}" github.com/c9s/fsrename/fsrename) \
        $(go list -f "{{.Root}}" github.com/c9s/fsrename/fsrename)/src/fsrename

    # see the document
    go doc fsrename



## OPTIONS

### FILTER OPTIONS

`-match [regexp pattern]` pre-filter the files and directories based on the given regular pattern.

`-contains [string]` pre-filter the files and directories based on the given string needle.

`-f`,`-file` only for files.

`-d`, `-dir` only for directories.

`-ext [extension name]` find files with matched file extension.

### REPLACEMENT OPTIONS

Please note the replacement target only works for the basename of a path.
`-replace*` and `-with*` should be combined together to replace the substrings.

#### Specifying replace target

`-r [search string]`, `-replace [search string]` specify target substring with normal string matching.

`-rr [regexp pattern]`, `-replace-re [regexp pattern]`, `-replace-regexp [regexp pattern]` specify target substring with regular expression matching.

#### Specifying replacement

`-w [string]`, `-with [string]` replacement for the target substring.

`-wf [format string]`, `-with-format [format string]` replacement with fmt.Sprintf format for the target substring.

### REPLACE RULE BUILDER OPTIONS

`-trim-prefix [prefix]` trim filename prefix.

`-trim-suffix [suffix]` trim filename suffix (this option removes suffix even for filename extensions).

`-camel` converts dash/underscore separated filenames into camelcase filenames.

`-underscore` converts camelcase filesnames into underscore separated filenames.

### Common Options

`-dryrun`  dry run, don't rename, just preview the result.

`-changelog [changelog file]` records the rename actions in CSV format file.



## QUICK EXAMPLES

Find files with extension `.php` and replace the substring from the filename.

    fsrename -ext "php" -replace "some" -with "others" src/

Replace `Stmt.go` with "_stmt.go" under the current directory:

    fsrename -replace "Stmt.go" -with "_stmt.go"

Replace `Stmt.go` with "_stmt.go" under directory `src/c6`:

    fsrename -replace "Stmt.go" -with "_stmt.go" src/c6

    # -r is a shorthand of -replace
    fsrename -r "Stmt.go" -with "_stmt.go" src/c6

Replace `foo` with `bar` from files contains `prefix_` 

    fsrename -file -contains prefix_ -replace foo -with bar test

Or use `-match` to pre-filter the files with regular expression

    fsrename -file -match "[a-z]+" -replace foo -with bar test

Use regular expression without escaping:

    fsrename -replace-regexp "_[a-z]*.go" -with ".go" src/c6

    # -rre is a shorthand of -replace-regexp
    fsrename -rre "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -replace-regexp "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -ext go -replace-regexp "[a-z]*" -with "123" src/c6

    fsrename -dir -replace "_xxx" -with "_aaa" src/c6

    fsrename -replace "_xxx" -with "_aaa" -dryrun  src/c6


## API

Build your own file stream workers.

The fsrename API is pretty straight forward and simple, you can create your own
filters, actors in just few lines:

```go
import "github.com/c9s/fsrename"


input := fsrename.NewFileStream()
scanner := fsrename.NewGlobScanner()
scanner.SetInput(input)
scanner.Start()
chain := scanner.
    Chain(fsrename.NewFileFilter()).
    Chain(fsrename.NewReverseSorter())

// start sending file entries into the input stream
input <- fsrename.MustNewFileEntry("tests")

// send EOS (end of stream)
input <- nil


// get entries from output
output := chain.Output()
entry := <-output
....
```


## Roadmap

v2.1

- [ ] Add `-list` to print the filtered file paths instead of renaming the files.
- [ ] Support more actions rather than file rename.
  - [ ] `-backup` to backup selected files.
- [ ] Add rename log printer to support rollback.
- [ ] Add `-grep` to grep files that contains specific pattern.
- [x] Ignore .git/.svn/.hg directories
- [ ] Add `-cleanup` to clean up non-ascii characters
- [ ] Add `--real` to solve symlink reference

## LICENSE

MIT License

