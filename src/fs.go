package fs

import (
	"fmt"
	"reflect"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/fatih/structs"
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

type FS struct {
	node fs.Node
}

type EntryGetter interface {
	GetDirentType() fuse.DirentType
}

var (
	inodeCnt   uint64
	Attributes fuse.Attr
	dataMap     map[string]any
)

const errNotPermitted = "Operation not permitted"

func Mount(mountPoint string) error {
	input := &structure{
		String: "str",
		Int:    18,
		Bool:   true,
		SubStructure: subStructure{
			Float: 1.3,
		},
	}

	c, err := fuse.Mount(mountPoint, fuse.FSName("fuse"), fuse.Subtype("tmpfs"))
	if err != nil {
		return err
	}
	defer c.Close()

	dataMap = structs.Map(input)

	err = fs.Serve(c, FS{})
	if err != nil {
		return err
	}

	return nil
}

func (FS) Root() (fs.Node, error) {
	dir := NewDir()
	dir.Entries = createEntries(dataMap)
	return dir, nil
}

func createEntries(structMap any) map[string]any {
	entries := map[string]any{}

	for key, value := range structMap.(map[string]any) {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			dir := NewDir()
			dir.Entries = createEntries(value)
			entries[key] = dir
		} else {
			entries[key] = NewFile([]byte(fmt.Sprint(reflect.ValueOf(value))))
		}
	}
	return entries
}
