Fast Rename
======================
A rename tool that supports concurrent renaming.

fsrename separates the pattern option, therefore you can specify the pattern without typeing these character escapes.

further, this tool supports concurrent renaming (powered by Go's channel and routine)

Install
--------------

    go get -u -x github.com/c9s/fsrename


Options
---------------

- `-match` the pattern that will match the files/dirs you want.
- `-replace` replace the matched string with what you want.
- `-dryrun`  dry run, don't rename, just preview the result.
- `-fileonly` rename only files.
- `-dironly` rename only directory.
- `-forext` rename only matched extension.
- `-c=2` number of workers. (concurrency)


Usage
---------------

    fsrename -match "_stmt.go" -replace "_stmt.go" src/c6

    fsrename -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileonly -forExt go -match "[a-z]*" -replace "123" src/c6

    fsrename -dironly -match "_xxx" -replace "_aaa" src/c6

    fsrename -dryrun -match "_xxx" -replace "_aaa" src/c6
