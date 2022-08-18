package fs

import (
	"fmt"
	"log"
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

func Mount(mountPoint string, userStruct any) error {
	c, err := fuse.Mount(mountPoint, fuse.FSName("fuse"), fuse.Subtype("tmpfs"))
	if err != nil {
		return err
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println("close: ", err)
		}
		fuse.Unmount(mountPoint)
	}()

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
	dir.Entries = f.createEntries(structMap)
	return dir, nil
}

func (f *FS) createEntries(structMap any) map[string]any {
	entries := map[string]any{}

	for key, value := range structMap.(map[string]any) {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			dir := NewDir()
			dir.Entries = f.createEntries(value)
			entries[key] = dir
		} else {
			content := []byte(fmt.Sprintln(reflect.ValueOf(value)))
			entries[key] = NewFile(key, f.userStruct, len(content))
		}
	}
	return entries
}
