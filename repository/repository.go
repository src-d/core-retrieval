// Package repository implements services to work with Git repository storage.
package repository

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/src-d/go-billy.v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// RootedTransactioner can initiate transactions on rooted repositories.
type RootedTransactioner interface {
	Begin(h plumbing.Hash) (Tx, error)
}

// Tx is a transaction on a repository. Any change performed in the given
// repository storer is in a transaction context. Transactions are guaranteed
// to be isolated.
type Tx interface {
	// Storer gets the repository storer. It returns the same instance on
	// every call until Commit or Rollback is performed.
	Storer() storage.Storer
	// Commit commits all changes to the repository.
	Commit() error
	// Rollback undoes any changes and cleans up.
	Rollback() error
}

type fsSrv struct {
	fs    billy.Filesystem
	local billy.Filesystem
}

// NewFilesystemRootedTransactioner returns a RootedTransactioner for repositories
// stored in the given billy.Filesystem, and uses a second billy.Filesystem
// as temporary storage for in-progress transactions.
//
// Note that transactionality is not fully guaranteed by this implementation,
// since it relies on recursive copying between arbitrary filesystems. If a
// Commit operation fails, the state of the first filesystem is unknown and can
// be invalid.
func NewFilesystemRootedTransactioner(fs, local billy.Filesystem) RootedTransactioner {
	return &fsSrv{fs, local}
}

func (s *fsSrv) Begin(h plumbing.Hash) (Tx, error) {
	origPath := h.String()
	tmpPath := fmt.Sprintf("%s/%d", h.String(), time.Now().UnixNano())
	if err := copyRecursive(s.fs, s.local, origPath, tmpPath); err != nil {
		return nil, err
	}

	sto, err := filesystem.NewStorage(s.local.Dir(tmpPath))
	if err != nil {
		return nil, err
	}

	if _, err := git.Open(sto, nil); err == git.ErrRepositoryNotExists {
		if _, err := git.Init(sto, nil); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &fsTx{
		fs:       s.fs,
		local:    s.local,
		origPath: origPath,
		tmpPath:  tmpPath,
		s:        sto,
	}, nil
}

type fsTx struct {
	fs, local         billy.Filesystem
	tmpPath, origPath string
	s                 storage.Storer
}

func (tx *fsTx) Storer() storage.Storer {
	return tx.s
}

func (tx *fsTx) Commit() error {
	if err := copyRecursive(tx.local, tx.fs, tx.tmpPath, tx.origPath); err != nil {
		_ = tx.cleanUp()
		return err
	}

	return tx.cleanUp()
}

func (tx *fsTx) Rollback() error {
	return tx.cleanUp()
}

func (tx *fsTx) cleanUp() error {
	return billy.RemoveAll(tx.local, tx.tmpPath)
}

func copyRecursive(fromFs, toFs billy.Filesystem, from, to string) (err error) {
	srcInfo, err := fromFs.Stat(from)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return copyFile(fromFs, toFs, from, to)
	}

	fis, err := fromFs.ReadDir(from)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		fromPath := fromFs.Join(from, fi.Name())
		toPath := toFs.Join(to, fi.Name())
		err := copyRecursive(fromFs, toFs, fromPath, toPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(fromFs, toFs billy.Filesystem, from, to string) (err error) {
	src, err := fromFs.Open(from)
	if err != nil {
		return err
	}
	defer checkClose(src, &err)

	dst, err := toFs.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer checkClose(dst, &err)

	_, err = io.Copy(dst, src)
	return err
}

func checkClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
