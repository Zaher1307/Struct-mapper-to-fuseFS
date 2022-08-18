package fs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type File struct {
	Type       fuse.DirentType
	Content    []byte
	Attributes fuse.Attr
}

var _ = (fs.Node)((*File)(nil))
var _ = (fs.HandleReadAller)((*File)(nil))
var _ = (fs.NodeSetattrer)((*File)(nil))
var _ = (EntryGetter)((*File)(nil))

func NewFile(content []byte) *File {
	return &File{
		Type:    fuse.DT_File,
		Content: content,
		Attributes: fuse.Attr{
			Inode: inodeCnt,
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
	fmt.Println(f.Content)
	return f.Content, nil
}

func (f *File) GetDirentType() fuse.DirentType {
	return f.Type
}

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	return errors.New(errNotPermitted)
}
