## package m

This package has convenient functions to create and play musical objects

    package main

    import (
        "github.com/emicklei/melrose/m"
        "fmt"
    )

    func main() {
        fmt.Println(m.Sequence("C# D_ E+"))
    }