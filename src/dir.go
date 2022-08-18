package fs

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type Dir struct {
	Type       fuse.DirentType
	Attributes fuse.Attr
	Entries    map[string]any
}

var _ = (fs.Node)((*Dir)(nil))
var _ = (fs.HandleReadDirAller)((*Dir)(nil))
var _ = (fs.NodeSetattrer)((*Dir)(nil))
var _ = (EntryGetter)((*Dir)(nil))

func NewDir() *Dir {
	return &Dir{
		Type: fuse.DT_Dir,
		Attributes: fuse.Attr{
			Inode: 0,
			Atime: time.Now(),
			Mtime: time.Now(),
			Ctime: time.Now(),
			Mode:  os.ModeDir | 0o555,
		},
		Entries: map[string]any{},
	}
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = d.Attributes
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	node, ok := d.Entries[name]
	if ok {
		return node.(fs.Node), nil
	}
	return nil, syscall.ENOENT
}

func (d *Dir) GetDirentType() fuse.DirentType {
	return d.Type
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries []fuse.Dirent

	for key, value := range d.Entries {
		var a fuse.Attr
		value.(fs.Node).Attr(ctx, &a)
		entries = append(entries, fuse.Dirent{
			Inode: a.Inode,
			Type:  value.(EntryGetter).GetDirentType(),
			Name:  key,
		})
	}
	return entries, nil
}

func (d *Dir) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	if req.Valid.Atime() {
		d.Attributes.Atime = req.Atime
	}
	if req.Valid.Mtime() {
		d.Attributes.Mtime = req.Mtime
	}
	if req.Valid.Size() {
		d.Attributes.Size = req.Size
	}
	return nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	return errors.New(errNotPermitted)
}
