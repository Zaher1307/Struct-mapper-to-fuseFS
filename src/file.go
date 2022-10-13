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
	Content    []byte
	FilePath   []string
	StructRef  any
}

var _ = (fs.Node)((*File)(nil))
var _ = (fs.HandleReadAller)((*File)(nil))
var _ = (fs.NodeSetattrer)((*File)(nil))
var _ = (EntryGetter)((*File)(nil))

// create new file
func newFile(fileName string, filePath []string, structRef any, contentSize int) *File {
	return &File{
		Type:      fuse.DT_File,
		FileName:  fileName,
		FilePath:  filePath,
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
	f.updateFileContent()
	*a = f.Attributes
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	f.updateFileContent()
	f.Attributes.Atime = time.Now()
	return f.Content, nil
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

// recursive function for fetching file content from struct reference
func (f *File) updateFileContent() {
	structMap := structs.Map(f.StructRef)

	for _, part := range f.FilePath {
		structMap = structMap[part].(map[string]any)
	}
	content := []byte(fmt.Sprintln(reflect.ValueOf(structMap[f.FileName])))

	f.Content = content
	f.Attributes.Size = uint64(len(content))
}
