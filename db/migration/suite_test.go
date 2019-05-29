package migration

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
	"testing"
)


func TestSuite_Apply(t *testing.T) {
	a := assert.New(t)
	testSchemaName := buildTestSchemaName()
	err := defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName))
	a.Nil(err)
	s := New(OptLog(logger.None()), OptGroups(createTestMigrations(testSchemaName)...))
	defer func() {
		// pq can't parameterize Drop
		err := defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName))
		a.Nil(err)
	}()
	err = s.Apply(context.Background(), defaultDB())
	a.Nil(err)

	ok, err := defaultDB().Query("SELECT 1 FROM pg_catalog.pg_indexes where indexname = $1 and tablename = $2", "idx_created_foo", "table_test_foo").Any()
	a.Nil(err)
	a.True(ok)

	ap, sk, fl, tot := s.Results()
	a.Equal(4, ap)
	a.Equal(1, sk)
	a.Equal(0, fl)
	a.Equal(5, tot)
}

func buildTestSchemaName() string {
	return fmt.Sprintf("test_sch_%s", stringutil.Random(stringutil.LowerLetters, 10))
}

func createTestMigrations(testSchemaName string) []*Group {
	return []*Group{
		NewGroupWithActions(
			NewStep(
				SchemaNotExists(testSchemaName),
				// pq can't parameterize Create
				func(i context.Context, connection *db.Connection, tx *sql.Tx) error {
					err := connection.Exec(fmt.Sprintf("CREATE SCHEMA %s;", testSchemaName))
					if err != nil {
						return err
					}
					(&connection.Config).Schema = testSchemaName
					return nil
				},
			)),
		NewGroupWithActions(
			NewStep(
				TableNotExists("table_test_foo"),
				Exec(fmt.Sprintf("CREATE TABLE %s.table_test_foo (id serial not null primary key, something varchar(32) not null);", testSchemaName)),
				),
			NewStep(
				ColumnNotExists("table_test_foo", "created_foo"),
				Exec(fmt.Sprintf("ALTER TABLE %s.table_test_foo ADD COLUMN created_foo timestamp not null;",testSchemaName)),
				)),
		NewGroup(OptSkipTransaction(), OptActions(
			NewStep(
				IndexNotExists("table_test_foo","idx_created_foo"),
				Exec(fmt.Sprintf("CREATE INDEX CONCURRENTLY idx_created_foo ON %s.table_test_foo(created_foo);", testSchemaName)),
				))),
		NewGroupWithActions(
			NewStep(
				TableNotExists("table_test_foo"),
				Exec(fmt.Sprintf("CREATE TABLE %s.table_test_foo (id serial not null primary key, something varchar(32) not null);", testSchemaName)),
			)),
	}
}