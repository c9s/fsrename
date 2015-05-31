FSRename
======================
A rename tool that supports concurrent renaming.

fsrename separates the pattern option, therefore you can specify the pattern without typeing these character escapes.

further, this tool supports concurrent renaming (powered by Go's channel and routine)

Install
--------------

    go get -u -x github.com/c9s/fsrename


Usage
---------------------

    fsrename [options] [path...]

When [path...] is not given, "./" will be used as the default path for scanning files.


Options
---------------

- `-match` the pattern that will match the files/dirs you want.
- `-replace` replace the matched string with what you want.

- `-trimprefix` trim filename prefix. When using this option, you don't have to specify `-match` or `-replace`.
- `-trimsuffix` trim filename suffix (including extension). When using this option, you don't have to specify `-match` or `-replace`.

- `-dryrun`  dry run, don't rename, just preview the result.
- `-fileonly` rename only files.
- `-dironly` rename only directory.
- `-forext` rename only matched extension.
- `-c=2` number of workers. (concurrency)


Examples
---------------

Replace `_stmt.go` with "_stmt.go" under the current directory:

    fsrename -match "_stmt.go" -replace "_stmt.go"

Replace `_stmt.go` with "_stmt.go" under directory `src/c6`:

    fsrename -match "_stmt.go" -replace "_stmt.go" src/c6

Use regular expression without escaping:

    fsrename -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -forExt go -match "[a-z]*" -replace "123" src/c6

    fsrename -dironly -match "_xxx" -replace "_aaa" src/c6

    fsrename -dryrun -match "_xxx" -replace "_aaa" src/c6
