package main

import (
	"fmt"
	"os"

	"fuse/src"
)

type structure struct {
	String       string
	Int          int
	Bool         bool
	SubStructure subStructure
}

type subStructure struct {
	Float float32
}

func main() {
	var err error
	if len(os.Args) != 2 {
		fmt.Println("too few arguments")
		fmt.Println(len(os.Args))
		os.Exit(1)
	}
	mountPoint := os.Args[1]
	input := &structure{
		String: "str",
		Int:    18,
		Bool:   true,
		SubStructure: subStructure{
			Float: 1.3,
		},
	}

	err = os.MkdirAll(mountPoint, 0555)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = fs.Mount(mountPoint, input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

