package core_retrieval

import (
	"database/sql"

	"srcd.works/core-retrieval.v0/model"

	"srcd.works/framework.v0/configurable"
	"srcd.works/framework.v0/database"
)

type containerConfig struct {
	configurable.BasicConfiguration
}

var config = &containerConfig{}

func init() {
	configurable.InitConfig(config)
}

var container struct {
	Database          *sql.DB
	ModelMentionStore *model.MentionStore
}

// Database returns a sql.DB for the default database. If it is not possible to
// connect to the database, this function will panic. Multiple calls will always
// return the same instance.
func Database() *sql.DB {
	if container.Database == nil {
		container.Database = database.Must(database.Default())
	}

	return container.Database
}

// ModelMentionStore returns the default *model.ModelMentionStore, using the
// default database. If it is not possible to connect to the database, this
// function will panic. Multiple calls will always return the same instance.
func ModelMentionStore() *model.MentionStore {
	if container.ModelMentionStore == nil {
		container.ModelMentionStore = model.NewMentionStore(Database())
	}

	return container.ModelMentionStore
}
