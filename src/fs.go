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

// mount your struct as a filesystem in the mount point
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

	err = fs.Serve(c, newFS(userStruct))
	if err != nil {
		return err
	}

	return nil
}

// creating new filesystem
func newFS(userStruct any) *FS {
	return &FS{
		userStruct: userStruct,
	}
}

func (f *FS) Root() (fs.Node, error) {
	dir := newDir()
	structMap := structs.Map(f.userStruct)
	dir.Entries = f.createEntries(structMap, []string{})
	return dir, nil
}

// recursive function to create file system entries(files, dirs) from user struct 
func (f *FS) createEntries(structMap map[string]any, currentPath []string) map[string]any {
	entries := map[string]any{}

	for key, val := range structMap {
		if reflect.TypeOf(val).Kind() == reflect.Map {
			dir := newDir()
			dir.Entries = f.createEntries(val.(map[string]any), append(currentPath, key))
			entries[key] = dir
		} else {
			filePath := make([]string, len(currentPath))
			copy(filePath, currentPath)
			content := []byte(fmt.Sprintln(reflect.ValueOf(val)))
			file := newFile(key, filePath, f.userStruct, len(content))
			entries[key] = file
		}
	}

	return entries
}
