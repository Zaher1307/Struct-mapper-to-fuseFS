package fs

import (
	"fmt"
	"reflect"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/fatih/structs"
)

type FS struct {
	userStruct any
}

type EntryGetter interface {
	GetDirentType() fuse.DirentType
}

const errNotPermitted = "Operation not permitted"

func Mount(mountPoint string, userStruct any) error {
	c, err := fuse.Mount(mountPoint, fuse.FSName("fuse"), fuse.Subtype("tmpfs"))
	if err != nil {
		return err
	}
	defer c.Close()

	err = fs.Serve(c, NewFS(userStruct))
	if err != nil {
		return err
	}

	return nil
}

func NewFS(userStruct any) *FS {
    return &FS{
    	userStruct: userStruct,
    }
}

func (f *FS) Root() (fs.Node, error) {
	dir := NewDir()
	structMap := structs.Map(f.userStruct)
	dir.Entries = createEntries(structMap)
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
