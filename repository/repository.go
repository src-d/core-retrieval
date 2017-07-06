// Package repository implements services to work with Git repository storage.
package repository

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"gopkg.in/src-d/go-billy-siva.v3"
	"gopkg.in/src-d/go-billy.v3"
	"gopkg.in/src-d/go-billy.v3/util"
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

// NewSivaRootedTransactioner returns a RootedTransactioner for repositories
// stored in the given billy.Filesystem (using siva file format), and uses a
// second billy.Filesystem as temporary storage for in-progress transactions.
//
// Note that transactionality is not fully guaranteed by this implementation,
// since it relies on copying between arbitrary filesystems. If a
// Commit operation fails, the state of the first filesystem is unknown and can
// be invalid.
func NewSivaRootedTransactioner(fs, local billy.Filesystem) RootedTransactioner {
	return &fsSrv{fs, local}
}

func (s *fsSrv) Begin(h plumbing.Hash) (Tx, error) {
	origPath := fmt.Sprintf("%s.siva", h)
	localPath := s.local.Join(h.String(), strconv.FormatInt(time.Now().UnixNano(), 10))
	localSivaPath := localPath+ ".siva"
	localTmpPath := localPath + ".tmp"

	if err := copyFile(s.fs, s.local, origPath, localSivaPath); err != nil {
		return nil, err
	}

	tmpFs, err := s.local.Chroot(localTmpPath)
	if err != nil {
		return nil, err
	}

	fs, err := sivafs.NewFilesystem(s.local, localSivaPath, tmpFs)
	if err != nil {
		return nil, err
	}

	sto, err := filesystem.NewStorage(fs)
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
		sivafs:   fs,
		origPath: origPath,
		tmpPath:  localSivaPath,
		s:        sto,
	}, nil
}

type fsTx struct {
	fs, local         billy.Filesystem
	sivafs            sivafs.SivaSync
	tmpPath, origPath string
	s                 storage.Storer
}

func (tx *fsTx) Storer() storage.Storer {
	return tx.s
}

func (tx *fsTx) Commit() error {
	if err := tx.sivafs.Sync(); err != nil {
		return err
	}

	if err := copyFile(tx.local, tx.fs, tx.tmpPath, tx.origPath); err != nil {
		_ = tx.cleanUp()
		return err
	}

	return tx.cleanUp()
}

func (tx *fsTx) Rollback() error {
	return tx.cleanUp()
}

func (tx *fsTx) cleanUp() error {
	return util.RemoveAll(tx.local, tx.tmpPath)
}

func copyFile(fromFs, toFs billy.Filesystem, from, to string) (err error) {
	src, err := fromFs.Open(from)
	if os.IsNotExist(err) {
		return nil
	}

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
