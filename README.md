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
- `-fileonly` filter files.
- `-dironly` filter directories.
- `-forext` find files with matched file extension.

## Replace Options

- `-replace` replace the matched string with what you want.
- `-with` replace the matched string with what you want.

## Helper Options (combines built-in rules)

- `-trimprefix` trim filename prefix. When using this option, you don't have to specify `-match` or `-replace`.
- `-trimsuffix` trim filename suffix (including extension). When using this option, you don't have to specify `-match` or `-replace`.

## Common Options

- `-dryrun`  dry run, don't rename, just preview the result.
- `-c=2` number of workers. (concurrency)

Some Examples
-------------

Replace `Stmt.go` with "_stmt.go" under the current directory:

    fsrename -match "Stmt.go" -replace "_stmt.go"

Replace `Stmt.go` with "_stmt.go" under directory `src/c6`:

    fsrename -match "Stmt.go" -replace "_stmt.go" src/c6

Use regular expression without escaping:

    fsrename -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -forExt go -match "[a-z]*" -replace "123" src/c6

    fsrename -dironly -match "_xxx" -replace "_aaa" src/c6

    fsrename -dryrun -match "_xxx" -replace "_aaa" src/c6
