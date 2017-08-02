package repository

import (
	"io"
	"os"
	"path"

	"github.com/colinmarc/hdfs"
	"gopkg.in/src-d/go-billy.v3"
)

// Copier is in charge either to obtain a file from the Remote filesystem implementation,
// or send it from local.
type Copier interface {
	CopyFromRemote(src, dst string, localFs billy.Filesystem) error
	CopyToRemote(src, dst string, localFs billy.Filesystem) error
}

// NewLocalCopier returns a Copier using as a remote a Billy filesystem
func NewLocalCopier(fs billy.Filesystem) Copier {
	return &LocalCopier{fs}
}

type LocalCopier struct {
	fs billy.Filesystem
}

func (c *LocalCopier) CopyFromRemote(src, dst string, localFs billy.Filesystem) error {
	return c.copyFile(c.fs, localFs, src, dst)
}

func (c *LocalCopier) CopyToRemote(src, dst string, localFs billy.Filesystem) error {
	return c.copyFile(localFs, c.fs, src, dst)
}

func (c *LocalCopier) copyFile(fromFs, toFs billy.Filesystem, from, to string) (err error) {
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

// NewHDFSCopier returns a copier using as a remote an HDFS cluster.
// URL is the hdfs connection URL and base is the base path to store all the files.
func NewHDFSCopier(URL string, base string) Copier {
	return &HDFSCopier{url: URL, base: base}
}

type HDFSCopier struct {
	url    string
	base   string
	client *hdfs.Client
}

// CopyFromRemote copies the file from HDFS to the provided billy Filesystem. If the file exists locally is overridden.
// If a writer is actually overriding the file on HDFS, we will able to read it, but a previous version of it.
func (c *HDFSCopier) CopyFromRemote(src, dst string, localFs billy.Filesystem) (err error) {
	if err := c.initializeClient(); err != nil {
		return err
	}

	rf, err := c.client.Open(path.Join(c.base, src))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer checkClose(rf, &err)

	lf, err := localFs.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer checkClose(lf, &err)

	_, err = io.Copy(lf, rf)
	return
}

// CopyToRemote copies from the provided billy Filesystem to HDFS. If the file exists on HDFS it will be overridden.
// If other writer is actually copying the same file to HDFS this method will throw an error because the WORM principle
// (Write Once Read Many).
func (c *HDFSCopier) CopyToRemote(src, dst string, localFs billy.Filesystem) (err error) {
	p := path.Join(c.base, dst)
	if err := c.initializeClient(); err != nil {
		return err
	}

	lf, err := localFs.Open(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer checkClose(lf, &err)

	if err := c.client.MkdirAll(path.Dir(p), os.FileMode(0644)); err != nil {
		return err
	}

	// TODO to avoid this, we should implement a 'truncate' flag in 'client.Create' method
	_, err = c.client.Stat(p)
	if err == nil {
		err = c.client.Remove(p)
		if err != nil {
			return err
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	rf, err := c.client.Create(p)
	if err != nil {
		return err
	}
	defer checkClose(rf, &err)

	_, err = io.Copy(rf, lf)
	return
}

func (c *HDFSCopier) initializeClient() (err error) {
	if c.client != nil {
		return
	}
	c.client, err = hdfs.New(c.url)

	return
}

func checkClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
