Fast Rename
======================

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
