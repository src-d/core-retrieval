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

func TestRootedTransactioner(t *testing.T) {
	require := require.New(t)

	fs := RootedTransactioner()
	require.NotNil(fs)

	fs2 := RootedTransactioner()
	require.Exactly(fs, fs2)
}
