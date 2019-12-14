package main

import "testing"
import "os"
import "bootstrap/asm"
import "bootstrap/cpu"
import "fmt"
import "io/ioutil"

func TestDummy(t *testing.T) {}

func TestS(t *testing.T) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		name := file.Name()

		if len(name) < 1 {
			continue
		}

		if name[0] == 'T' {
			fmt.Printf("Running %s\n", name)
			run(name)
		}
	}
}

func run(fpath string) {
	f, err := os.Open(fpath)

	if err != nil {
		panic(err.Error())
	}

	regs := make([]uint32, 0x10)
	ec := asm.NewEmitContext(1024, 4096)
	asm.AsmReader(ec, f)
	ec.Resolve()

	if cpu.Execute(regs, ec.Memory()) != 0 {
		panic(fpath + " has failed!")
	}
}
