package core_retrieval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	require := require.New(t)
	db := Database()
	require.NotNil(db)

	db2 := Database()
	require.Exactly(db, db2)
}

func TestModelMentionStore(t *testing.T) {
	require := require.New(t)
	s := ModelMentionStore()
	require.NotNil(s)

	s2 := ModelMentionStore()
	require.Exactly(s, s2)
}

func TestFilesystemRootedTransactioner(t *testing.T) {
	require := require.New(t)

	fs := FilesystemRootedTransactioner()
	require.NotNil(fs)

	fs2 := FilesystemRootedTransactioner()
	require.Exactly(fs, fs2)
}
