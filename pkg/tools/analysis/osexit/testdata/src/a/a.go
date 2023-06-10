package main

import "os"

func main() {
	os.Exit(0) // want "Found call of os.Exit on main package"
}
