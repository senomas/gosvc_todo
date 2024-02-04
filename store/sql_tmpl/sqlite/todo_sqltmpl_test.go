package sqlite_test

import (
	"testing"

	svc_store "github.com/senomas/gosvc_store/store"
	"github.com/senomas/gosvc_todo/store"
	"github.com/senomas/gosvc_todo/store/sql_tmpl"
	"github.com/senomas/gosvc_todo/store/sql_tmpl/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestTodoSqlTemplateSqlite(t *testing.T) {
	var tl sql_tmpl.TodoStoreTemplate = &sqlite.TodoStoreTemplateImpl{}
	assert.NotNil(t, tl, "TodoStoreTemplateImpl not nil")

	t.Run("FindTodoTotal no filter", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{})
		assert.Equal(t, `SELECT COUNT(*) FROM todo`, qry)
		assert.EqualValues(t, len(args), 0, "no args")
	})

	t.Run("FindTodoTotal where id =", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{
			ID: svc_store.FilterInt64{Value: 100, Op: svc_store.OP_EQ},
		})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE id = $1`, qry)
		assert.EqualValues(t, args, []any{int64(100)}, "1 args")
	})

	t.Run("FindTodoTotal where title =", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{
			Title: svc_store.FilterString{Value: "foo", Op: svc_store.OP_EQ},
		})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE title = $1`, qry)
		assert.EqualValues(t, args, []any{"foo"}, "1 args")
	})

	t.Run("FindTodoTotal where title like", func(t *testing.T) {
		qry, args := tl.FindTodoTotal(store.TodoFilter{
			Title: svc_store.FilterString{Value: "foo", Op: svc_store.OP_LIKE},
		})
		assert.Equal(t, `SELECT COUNT(*) FROM todo WHERE title LIKE $1`, qry)
		assert.EqualValues(t, args, []any{"foo"}, "1 args")
	})
}
