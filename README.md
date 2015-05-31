Fast Rename
======================
A rename tool that supports concurrent renaming.

File::Rename is written in Perl, it works fine. however when you need to pass
complex regular expression, you have to escape the characters, e.g. the slash `/`.

fsrename separate the pattern option, therefore you can specify the pattern without typeing these character escapes.

further, this tool supports concurrent renaming (powered by Go's channel and routine)

Install
--------------

    go get -u -x github.com/c9s/fsrename


Usage
---------------


    fsrename -match "_stmt.go" -replace "_stmt.go" src/c6

    fsrename -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileOnly -match "_[a-z]*.go" -replace ".go" src/c6

    fsrename -fileOnly -forExt go -match "[a-z]*" -replace "123" src/c6

    fsrename -dirOnly -match "_xxx" -replace "_aaa" src/c6
