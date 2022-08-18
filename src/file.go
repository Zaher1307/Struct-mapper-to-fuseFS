package fs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/fatih/structs"
)

type File struct {
	Type       fuse.DirentType
	Attributes fuse.Attr
	FileName   string
	StructRef  any
}

var _ = (fs.Node)((*File)(nil))
var _ = (fs.HandleReadAller)((*File)(nil))
var _ = (fs.NodeSetattrer)((*File)(nil))
var _ = (EntryGetter)((*File)(nil))

func NewFile(fileName string, structRef any, contentSize int) *File {
	return &File{
		Type:    fuse.DT_File,
		FileName: fileName,
		StructRef: structRef,
		Attributes: fuse.Attr{
			Inode: 0,
			Size:  uint64(contentSize),
			Atime: time.Now(),
			Mtime: time.Now(),
			Ctime: time.Now(),
			Mode:  0o444,
		},
	}
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = f.Attributes
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return f.fetchFileContent(), nil
}

func (f *File) GetDirentType() fuse.DirentType {
	return f.Type
}

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	if req.Valid.Atime() {
		f.Attributes.Atime = req.Atime
	}

	if req.Valid.Mtime() {
		f.Attributes.Mtime = req.Mtime
	}

	if req.Valid.Size() {
		f.Attributes.Size = req.Size
	}

	return nil
}

func (f *File) fetchFileContent() []byte {
    structMap := structs.Map(f.StructRef)
    var result []byte
    var traverse func(map[string]any)

    traverse = func(m map[string]any) {
        for key, val := range m {
            if reflect.TypeOf(val).Kind() == reflect.Map {
                traverse(val.(map[string]any))
            } else {
                if key == f.FileName {
                    result = []byte(fmt.Sprintln(reflect.ValueOf(val)))
                    f.Attributes.Size = uint64(len(result))
                }
            }
        }
    }

    traverse(structMap)

    return result
}
