package main

import (
    "fmt"
    "github.com/KoFish/pallium/rest"
    "github.com/KoFish/pallium/storage"
)

func main() {
    fmt.Println("Setting up matrix")
    storage.Setup()
    rest.Setup()
    fmt.Println("Starting service")
    rest.Start()
    fmt.Println("Shutting down matrix")
}
