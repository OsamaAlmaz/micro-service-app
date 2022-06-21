package main

import (
    "fmt"
)

var a int 

func main() {
    a = 1
    readfile()
}

func readfile(){
    fmt.Println("This is a", a)
}
