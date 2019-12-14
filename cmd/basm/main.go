package main

import "bootstrap/asm"

import "fmt"
import "flag"
import "os"
import "bufio"
import "strings"

func main() {
	fmt.Println("basm - bootstrap assembler\n")

	mod_file := flag.String("module", "", "Path to module file (*.M)")
	importPaths := flag.String("paths", "", "Import paths")
	out := flag.String("out", "output.bin", "Output path")

	flag.Parse()

	if *mod_file == "" {
		fmt.Fprintf(os.Stderr, "No module file specified!")
		os.Exit(1)
	}

	ec := asm.NewDefaultEmitContext()
	processModuleFile(ec, *mod_file, strings.Split(*importPaths, ":"))

	xunresolved := ec.GetXUnresolved()

	if len(xunresolved) == 0 {
		fmt.Fprintf(os.Stdout, "Writing to output file\n")

		f, err := os.Create(*out)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create output file %q: %q!", *out, err.Error())
			os.Exit(1)
		}

		_, err = f.Write(ec.Memory())

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not write to output file %q: %q!", *out, err.Error())
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Success\n")
		os.Exit(0)
	} else {
		for _, name := range xunresolved {
			fmt.Fprintf(os.Stderr, "Unresolved %q!\n", name)
		}
		os.Exit(1)
	}
}

func openFile(path string, importPaths []string) (*os.File, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		for _, importPath := range importPaths {
			_, err := os.Stat(importPath + "/" + path)

			if os.IsNotExist(err) {
				continue
			}

			if err == nil {
				fmt.Fprintf(os.Stdout, "Found %q in %q\n", path, importPath)

				return os.Open(importPath + "/" + path)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "Found %q in \".\"\n", path)

	return os.Open(path)
}

func processModuleFile(ec *asm.EmitContext, mod_file string, importPaths []string) {
	f, err := openFile(mod_file, importPaths)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open module file %q: %q!", mod_file, err.Error())
		os.Exit(1)
	}

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		path := sc.Text()

		if strings.HasSuffix(path, ".M") {
			processModuleFile(ec, path, importPaths)
		} else if strings.HasSuffix(path, ".S") {
			fh, err := openFile(path, importPaths)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not open assembly file %q: %q!", path, err.Error())
				os.Exit(1)
			}

			asm.AsmReader(ec, fh)
			ec.Resolve()
		}
	}
}
