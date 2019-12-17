package main

import "bootstrap/cpu"
import "flag"
import "os"
import "io/ioutil"
import "fmt"

func main() {
	fileptr := flag.String("bin", "", "A .bin file to be run")

	flag.Parse()

	if *fileptr == "" {
		panic("No file specified!")
	}

	f, err := os.Open(*fileptr)

	if err != nil {
		panic(err.Error())
	}
	
	fileContents, err := ioutil.ReadAll(f)
	
	if err != nil {
		panic(err.Error())
	}
	
	memory := make([]uint8, 1024 + len(fileContents))
	copy(memory[1024:], fileContents)

	regs := make([]uint32, 0x10)

	if cpu.Execute(regs, memory) != 0 {
		fmt.Printf("Error on execute!")
		os.Exit(1)
	}
	
	os.Exit(0)
}
