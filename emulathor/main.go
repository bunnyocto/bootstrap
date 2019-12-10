package main

import "bootstrap/asm"
import "bootstrap/cpu"
import "flag"
import "os"

func main() {
	sfileptr := flag.String("sfile","","A *.S file to be run!")
	
	flag.Parse()
	
	if *sfileptr == "" {
		panic("No `sfile` specified!")
	}
	
	f, err := os.Open(*sfileptr)
	
	if err != nil {
		panic(err.Error())
	}
	
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()
	asm.AsmReader(ec, f)
	ec.Resolve()
	
	cpu.Execute(regs, ec.Memory())
}