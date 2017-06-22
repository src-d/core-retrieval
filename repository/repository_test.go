package repository

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/go-billy.v3"
	"gopkg.in/src-d/go-billy.v3/memfs"
	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	tmpPrefix = "core-retrieval-test"
)

var (
	h1 = plumbing.NewHash("0000000000000000000000000000000000000001")
)

func TestFilesystemSuite(t *testing.T) {
	suite.Run(t, &FilesystemSuite{})
}

type FilesystemSuite struct {
	suite.Suite
	tmpDirs map[string]bool
}

func (s *FilesystemSuite) SetupTest() {
	s.tmpDirs = make(map[string]bool)
}

func (s *FilesystemSuite) TearDownTest() {
	s.cleanUpTempDirectories()
}

func (s *FilesystemSuite) cleanUpTempDirectories() {
	require := require.New(s.T())
	var err error
	for dir := range s.tmpDirs {
		if e := os.RemoveAll(dir); e != nil && err != nil {
			err = e
		}

		delete(s.tmpDirs, dir)
	}

	require.NoError(err)
}

func (s *FilesystemSuite) Test() {
	fsPairs := []*fsPair{
		{"mem to mem", memfs.New(), memfs.New()},
		{"mem to os", memfs.New(), s.newFilesystem()},
		{"os to mem", s.newFilesystem(), memfs.New()},
		{"os to os", s.newFilesystem(), s.newFilesystem()},
	}

	for _, fsPair := range fsPairs {
		s.T().Run(fsPair.Name, func(t *testing.T) {
			testRootedTransactioner(t, NewSivaRootedTransactioner(fsPair.From, fsPair.To))
		})
	}
}

func (s *FilesystemSuite) newFilesystem() billy.Filesystem {
	require := require.New(s.T())
	tmpDir, err := ioutil.TempDir(os.TempDir(), tmpPrefix)
	require.NoError(err)
	s.tmpDirs[tmpDir] = true
	return osfs.New(tmpDir)
}

func testRootedTransactioner(t *testing.T, s RootedTransactioner) {
	require := require.New(t)

	// begin tx1
	tx1, err := s.Begin(h1)
	require.NoError(err)
	require.NotNil(tx1)
	r1, err := git.Open(tx1.Storer(), nil)
	require.NoError(err)

	// tx1 -> create ref1
	refName1 := plumbing.ReferenceName("ref1")
	err = r1.Storer.SetReference(plumbing.NewSymbolicReference(refName1, refName1))
	require.NoError(err)

	// begin tx2
	tx2, err := s.Begin(h1)
	require.NoError(err)
	require.NotNil(tx2)
	r2, err := git.Open(tx2.Storer(), nil)
	require.NoError(err)

	// ref1 not visible in tx2
	_, err = r2.Reference(refName1, false)
	require.Equal(plumbing.ErrReferenceNotFound, err)

	// tx2 -> create ref2
	refName2 := plumbing.ReferenceName("ref2")
	err = r2.Storer.SetReference(plumbing.NewSymbolicReference(refName2, refName2))
	require.NoError(err)

	// ref2 not visible in tx2
	_, err = r1.Reference(refName2, false)
	require.Equal(plumbing.ErrReferenceNotFound, err)

	// commit tx1
	err = tx1.Commit()
	require.NoError(err)

	// ref1 not visible in tx2 (even with tx1 committed)
	_, err = r2.Reference(refName1, false)
	require.Equal(plumbing.ErrReferenceNotFound, err)

	// rollback tx2
	err = tx2.Rollback()
	require.NoError(err)

	// begin tx3
	tx3, err := s.Begin(h1)
	require.NoError(err)
	require.NotNil(tx3)
	r3, err := git.Open(tx3.Storer(), nil)
	require.NoError(err)

	// ref1 visible in tx3
	_, err = r3.Reference(refName1, false)
	require.NoError(err)
	require.NoError(tx3.Rollback())
}

type fsPair struct {
	Name string
	From billy.Filesystem
	To   billy.Filesystem
}
