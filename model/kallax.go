// IMPORTANT! This is auto generated code by https://github.com/src-d/go-kallax
// Please, do not touch the code below, and if you do, do it under your own
// risk. Take into account that all the code you write here will be completely
// erased from earth the next time you generate the kallax models.
package model

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/src-d/go-kallax.v1"
	"gopkg.in/src-d/go-kallax.v1/types"
)

var _ types.SQLType
var _ fmt.Formatter

// NewMention returns a new instance of Mention.
func NewMention() (record *Mention) {
	return newMention()
}

// GetID returns the primary key of the model.
func (r *Mention) GetID() kallax.Identifier {
	return (*kallax.ULID)(&r.ID)
}

// ColumnAddress returns the pointer to the value of the given column.
func (r *Mention) ColumnAddress(col string) (interface{}, error) {
	switch col {
	case "id":
		return (*kallax.ULID)(&r.ID), nil
	case "created_at":
		return &r.Timestamps.CreatedAt, nil
	case "updated_at":
		return &r.Timestamps.UpdatedAt, nil
	case "endpoint":
		return &r.Endpoint, nil
	case "provider":
		return &r.Provider, nil
	case "vcs":
		return &r.VCS, nil
	case "context":
		return types.JSON(&r.Context), nil

	default:
		return nil, fmt.Errorf("kallax: invalid column in Mention: %s", col)
	}
}

// Value returns the value of the given column.
func (r *Mention) Value(col string) (interface{}, error) {
	switch col {
	case "id":
		return r.ID, nil
	case "created_at":
		return r.Timestamps.CreatedAt, nil
	case "updated_at":
		return r.Timestamps.UpdatedAt, nil
	case "endpoint":
		return r.Endpoint, nil
	case "provider":
		return r.Provider, nil
	case "vcs":
		return (string)(r.VCS), nil
	case "context":
		return types.JSON(r.Context), nil

	default:
		return nil, fmt.Errorf("kallax: invalid column in Mention: %s", col)
	}
}

// NewRelationshipRecord returns a new record for the relatiobship in the given
// field.
func (r *Mention) NewRelationshipRecord(field string) (kallax.Record, error) {
	return nil, fmt.Errorf("kallax: model Mention has no relationships")
}

// SetRelationship sets the given relationship in the given field.
func (r *Mention) SetRelationship(field string, rel interface{}) error {
	return fmt.Errorf("kallax: model Mention has no relationships")
}

// MentionStore is the entity to access the records of the type Mention
// in the database.
type MentionStore struct {
	*kallax.Store
}

// NewMentionStore creates a new instance of MentionStore
// using a SQL database.
func NewMentionStore(db *sql.DB) *MentionStore {
	return &MentionStore{kallax.NewStore(db)}
}

// Insert inserts a Mention in the database. A non-persisted object is
// required for this operation.
func (s *MentionStore) Insert(record *Mention) error {

	if err := record.BeforeSave(); err != nil {
		return err
	}

	return s.Store.Insert(Schema.Mention.BaseSchema, record)

}

// Update updates the given record on the database. If the columns are given,
// only these columns will be updated. Otherwise all of them will be.
// Be very careful with this, as you will have a potentially different object
// in memory but not on the database.
// Only writable records can be updated. Writable objects are those that have
// been just inserted or retrieved using a query with no custom select fields.
func (s *MentionStore) Update(record *Mention, cols ...kallax.SchemaField) (updated int64, err error) {

	if err := record.BeforeSave(); err != nil {
		return 0, err
	}

	return s.Store.Update(Schema.Mention.BaseSchema, record, cols...)

}

// Save inserts the object if the record is not persisted, otherwise it updates
// it. Same rules of Update and Insert apply depending on the case.
func (s *MentionStore) Save(record *Mention) (updated bool, err error) {
	if !record.IsPersisted() {
		return false, s.Insert(record)
	}

	rowsUpdated, err := s.Update(record)
	if err != nil {
		return false, err
	}

	return rowsUpdated > 0, nil
}

// Delete removes the given record from the database.
func (s *MentionStore) Delete(record *Mention) error {

	return s.Store.Delete(Schema.Mention.BaseSchema, record)

}

// Find returns the set of results for the given query.
func (s *MentionStore) Find(q *MentionQuery) (*MentionResultSet, error) {
	rs, err := s.Store.Find(q)
	if err != nil {
		return nil, err
	}

	return NewMentionResultSet(rs), nil
}

// MustFind returns the set of results for the given query, but panics if there
// is any error.
func (s *MentionStore) MustFind(q *MentionQuery) *MentionResultSet {
	return NewMentionResultSet(s.Store.MustFind(q))
}

// Count returns the number of rows that would be retrieved with the given
// query.
func (s *MentionStore) Count(q *MentionQuery) (int64, error) {
	return s.Store.Count(q)
}

// MustCount returns the number of rows that would be retrieved with the given
// query, but panics if there is an error.
func (s *MentionStore) MustCount(q *MentionQuery) int64 {
	return s.Store.MustCount(q)
}

// FindOne returns the first row returned by the given query.
// `ErrNotFound` is returned if there are no results.
func (s *MentionStore) FindOne(q *MentionQuery) (*Mention, error) {
	q.Limit(1)
	q.Offset(0)
	rs, err := s.Find(q)
	if err != nil {
		return nil, err
	}

	if !rs.Next() {
		return nil, kallax.ErrNotFound
	}

	record, err := rs.Get()
	if err != nil {
		return nil, err
	}

	if err := rs.Close(); err != nil {
		return nil, err
	}

	return record, nil
}

// MustFindOne returns the first row retrieved by the given query. It panics
// if there is an error or if there are no rows.
func (s *MentionStore) MustFindOne(q *MentionQuery) *Mention {
	record, err := s.FindOne(q)
	if err != nil {
		panic(err)
	}
	return record
}

// Reload refreshes the Mention with the data in the database and
// makes it writable.
func (s *MentionStore) Reload(record *Mention) error {
	return s.Store.Reload(Schema.Mention.BaseSchema, record)
}

// Transaction executes the given callback in a transaction and rollbacks if
// an error is returned.
// The transaction is only open in the store passed as a parameter to the
// callback.
func (s *MentionStore) Transaction(callback func(*MentionStore) error) error {
	if callback == nil {
		return kallax.ErrInvalidTxCallback
	}

	return s.Store.Transaction(func(store *kallax.Store) error {
		return callback(&MentionStore{store})
	})
}

// MentionQuery is the object used to create queries for the Mention
// entity.
type MentionQuery struct {
	*kallax.BaseQuery
}

// NewMentionQuery returns a new instance of MentionQuery.
func NewMentionQuery() *MentionQuery {
	return &MentionQuery{
		BaseQuery: kallax.NewBaseQuery(Schema.Mention.BaseSchema),
	}
}

// Select adds columns to select in the query.
func (q *MentionQuery) Select(columns ...kallax.SchemaField) *MentionQuery {
	if len(columns) == 0 {
		return q
	}
	q.BaseQuery.Select(columns...)
	return q
}

// SelectNot excludes columns from being selected in the query.
func (q *MentionQuery) SelectNot(columns ...kallax.SchemaField) *MentionQuery {
	q.BaseQuery.SelectNot(columns...)
	return q
}

// Copy returns a new identical copy of the query. Remember queries are mutable
// so make a copy any time you need to reuse them.
func (q *MentionQuery) Copy() *MentionQuery {
	return &MentionQuery{
		BaseQuery: q.BaseQuery.Copy(),
	}
}

// Order adds order clauses to the query for the given columns.
func (q *MentionQuery) Order(cols ...kallax.ColumnOrder) *MentionQuery {
	q.BaseQuery.Order(cols...)
	return q
}

// BatchSize sets the number of items to fetch per batch when there are 1:N
// relationships selected in the query.
func (q *MentionQuery) BatchSize(size uint64) *MentionQuery {
	q.BaseQuery.BatchSize(size)
	return q
}

// Limit sets the max number of items to retrieve.
func (q *MentionQuery) Limit(n uint64) *MentionQuery {
	q.BaseQuery.Limit(n)
	return q
}

// Offset sets the number of items to skip from the result set of items.
func (q *MentionQuery) Offset(n uint64) *MentionQuery {
	q.BaseQuery.Offset(n)
	return q
}

// Where adds a condition to the query. All conditions added are concatenated
// using a logical AND.
func (q *MentionQuery) Where(cond kallax.Condition) *MentionQuery {
	q.BaseQuery.Where(cond)
	return q
}

// FindByID adds a new filter to the query that will require that
// the ID property is equal to one of the passed values; if no passed values, it will do nothing
func (q *MentionQuery) FindByID(v ...kallax.ULID) *MentionQuery {
	if len(v) == 0 {
		return q
	}
	values := make([]interface{}, len(v))
	for i, val := range v {
		values[i] = val
	}
	return q.Where(kallax.In(Schema.Mention.ID, values...))
}

// FindByCreatedAt adds a new filter to the query that will require that
// the CreatedAt property is equal to the passed value
func (q *MentionQuery) FindByCreatedAt(cond kallax.ScalarCond, v time.Time) *MentionQuery {
	return q.Where(cond(Schema.Mention.CreatedAt, v))
}

// FindByUpdatedAt adds a new filter to the query that will require that
// the UpdatedAt property is equal to the passed value
func (q *MentionQuery) FindByUpdatedAt(cond kallax.ScalarCond, v time.Time) *MentionQuery {
	return q.Where(cond(Schema.Mention.UpdatedAt, v))
}

// FindByEndpoint adds a new filter to the query that will require that
// the Endpoint property is equal to the passed value
func (q *MentionQuery) FindByEndpoint(v string) *MentionQuery {
	return q.Where(kallax.Eq(Schema.Mention.Endpoint, v))
}

// FindByProvider adds a new filter to the query that will require that
// the Provider property is equal to the passed value
func (q *MentionQuery) FindByProvider(v string) *MentionQuery {
	return q.Where(kallax.Eq(Schema.Mention.Provider, v))
}

// FindByVCS adds a new filter to the query that will require that
// the VCS property is equal to the passed value
func (q *MentionQuery) FindByVCS(v VCS) *MentionQuery {
	return q.Where(kallax.Eq(Schema.Mention.VCS, v))
}

// MentionResultSet is the set of results returned by a query to the
// database.
type MentionResultSet struct {
	ResultSet kallax.ResultSet
	last      *Mention
	lastErr   error
}

// NewMentionResultSet creates a new result set for rows of the type
// Mention.
func NewMentionResultSet(rs kallax.ResultSet) *MentionResultSet {
	return &MentionResultSet{ResultSet: rs}
}

// Next fetches the next item in the result set and returns true if there is
// a next item.
// The result set is closed automatically when there are no more items.
func (rs *MentionResultSet) Next() bool {
	if !rs.ResultSet.Next() {
		rs.lastErr = rs.ResultSet.Close()
		rs.last = nil
		return false
	}

	var record kallax.Record
	record, rs.lastErr = rs.ResultSet.Get(Schema.Mention.BaseSchema)
	if rs.lastErr != nil {
		rs.last = nil
	} else {
		var ok bool
		rs.last, ok = record.(*Mention)
		if !ok {
			rs.lastErr = fmt.Errorf("kallax: unable to convert record to *Mention")
			rs.last = nil
		}
	}

	return true
}

// Get retrieves the last fetched item from the result set and the last error.
func (rs *MentionResultSet) Get() (*Mention, error) {
	return rs.last, rs.lastErr
}

// ForEach iterates over the complete result set passing every record found to
// the given callback. It is possible to stop the iteration by returning
// `kallax.ErrStop` in the callback.
// Result set is always closed at the end.
func (rs *MentionResultSet) ForEach(fn func(*Mention) error) error {
	for rs.Next() {
		record, err := rs.Get()
		if err != nil {
			return err
		}

		if err := fn(record); err != nil {
			if err == kallax.ErrStop {
				return rs.Close()
			}

			return err
		}
	}
	return nil
}

// All returns all records on the result set and closes the result set.
func (rs *MentionResultSet) All() ([]*Mention, error) {
	var result []*Mention
	for rs.Next() {
		record, err := rs.Get()
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// One returns the first record on the result set and closes the result set.
func (rs *MentionResultSet) One() (*Mention, error) {
	if !rs.Next() {
		return nil, kallax.ErrNotFound
	}

	record, err := rs.Get()
	if err != nil {
		return nil, err
	}

	if err := rs.Close(); err != nil {
		return nil, err
	}

	return record, nil
}

// Err returns the last error occurred.
func (rs *MentionResultSet) Err() error {
	return rs.lastErr
}

// Close closes the result set.
func (rs *MentionResultSet) Close() error {
	return rs.ResultSet.Close()
}

type schema struct {
	Mention *schemaMention
}

type schemaMention struct {
	*kallax.BaseSchema
	ID        kallax.SchemaField
	CreatedAt kallax.SchemaField
	UpdatedAt kallax.SchemaField
	Endpoint  kallax.SchemaField
	Provider  kallax.SchemaField
	VCS       kallax.SchemaField
	Context   kallax.SchemaField
}

var Schema = &schema{
	Mention: &schemaMention{
		BaseSchema: kallax.NewBaseSchema(
			"mentions",
			"__mention",
			kallax.NewSchemaField("id"),
			kallax.ForeignKeys{},
			func() kallax.Record {
				return new(Mention)
			},
			false,
			kallax.NewSchemaField("id"),
			kallax.NewSchemaField("created_at"),
			kallax.NewSchemaField("updated_at"),
			kallax.NewSchemaField("endpoint"),
			kallax.NewSchemaField("provider"),
			kallax.NewSchemaField("vcs"),
			kallax.NewSchemaField("context"),
		),
		ID:        kallax.NewSchemaField("id"),
		CreatedAt: kallax.NewSchemaField("created_at"),
		UpdatedAt: kallax.NewSchemaField("updated_at"),
		Endpoint:  kallax.NewSchemaField("endpoint"),
		Provider:  kallax.NewSchemaField("provider"),
		VCS:       kallax.NewSchemaField("vcs"),
		Context:   kallax.NewSchemaField("context"),
	},
}
