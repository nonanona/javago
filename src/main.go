package main

import (
  "os"
  "vm"
)

func main() {
  vm := vm.NewVM()
  vm.Execute(os.Args[1])
}
