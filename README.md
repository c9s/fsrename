FSRename
======================
A simple, powerful rename tool supports complex filtering

fsrename separates the pattern option, therefore you can specify the pattern without typeing these character escapes.

further, this tool supports concurrent renaming (powered by Go's channel and routine)

Install
--------------

    go get -u -x github.com/c9s/fsrename


Usage
---------------------

    fsrename [options] [path...]

When [path...] is not given, "./" will be used as the default path for scanning files.

## Filter Options

- `-match` match option filters the files and directories based on the given regular pattern.
- `-file` only for files.
- `-dir` only for directories.
- `-forext` find files with matched file extension.

## Replace Options

Please note the replacement target only works for the basename of a path.

- `-replace` specify target substring with normal string matching.
- `-replaceRegexp` specify target substring with regular expression matching.

- `-with` replacement for the target substring.
- `-withFormat` replacement with fmt.Sprintf format for the target substring.

## Replace Rule Builder Options

- `-trimPrefix` trim filename prefix.
- `-trimSuffix` trim filename suffix (this option removes suffix even for filename extensions).

## Common Options

- `-dryrun`  dry run, don't rename, just preview the result.

Quick Examples
-------------

Replace `Stmt.go` with "_stmt.go" under the current directory:

    fsrename -replace "Stmt.go" -with "_stmt.go"

Replace `Stmt.go` with "_stmt.go" under directory `src/c6`:

    fsrename -replace "Stmt.go" -with "_stmt.go" src/c6

Replace `foo` with `bar` from files contains `prefix_` 

    fsrename -file -match prefix_ -replace foo -with bar test

Use regular expression without escaping:

    fsrename -replaceRegexp "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -replaceRegexp "_[a-z]*.go" -with ".go" src/c6

    fsrename -file -ext go -replaceRegexp "[a-z]*" -with "123" src/c6

    fsrename -dir -replace "_xxx" -with "_aaa" src/c6

    fsrename -replace "_xxx" -with "_aaa" -dryrun  src/c6
