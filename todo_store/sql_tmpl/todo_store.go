package sql_tmpl

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/senomas/gosvc_store/store"
	"github.com/senomas/gosvc_todo/todo_store"
)

var (
	errCtxNoDB          = errors.New("no db defined in context")
	errNotodo_storeTmpl = errors.New("todo todo_store template not initialized")
)

type TodoStoreImpl struct{}

type TodoStoreTemplate interface {
	InsertTodo(t *todo_store.Todo) (string, []any)
	UpdateTodo(t *todo_store.Todo) (string, []any)
	DeleteTodoByID(id any) (string, []any)

	GetTodoByID(id any) (string, []any)
	FindTodo(todo_store.TodoFilter, int64, int) (string, []any)
	FindTodoTotal(todo_store.TodoFilter) (string, []any)

	ErrorMapFind(error) error
}

func init() {
	slog.Debug("Register sql_tmpl.Todotodo_store")
	todo_store.SetupTodoStoreImplementation(&TodoStoreImpl{})
}

var todotodo_storeTemplateImpl TodoStoreTemplate

func SetupTodoStoreTemplate(t TodoStoreTemplate) {
	todotodo_storeTemplateImpl = t
}

func (t *TodoStoreImpl) Init(ctx context.Context) error {
	if todotodo_storeTemplateImpl == nil {
		return errNotodo_storeTmpl
	}
	return nil
}

// CreateTodo implements todo_store.Todotodo_store.
func (t *TodoStoreImpl) CreateTodo(ctx context.Context, title string) (*todo_store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		todo := todo_store.Todo{Title: title}
		qry, args := todotodo_storeTemplateImpl.InsertTodo(&todo)
		slog.Debug("CreateTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		rs, err := db.ExecContext(ctx, qry, args...)
		if err != nil {
			slog.Warn("Error insert todo", "qry", qry, "error", err)
			return nil, err
		}
		todo.ID, err = rs.LastInsertId()
		return &todo, err
	}
	return nil, errCtxNoDB
}

// UpdateTodo implements todo_store.Todotodo_store.
func (t *TodoStoreImpl) UpdateTodo(ctx context.Context, todo todo_store.Todo) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := todotodo_storeTemplateImpl.UpdateTodo(&todo)
		slog.Debug("UpdateTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		_, err := db.ExecContext(ctx, qry, args...)
		return err
	}
	return errCtxNoDB
}

// DeleteTodoByID implements todo_store.Todotodo_store.
func (t *TodoStoreImpl) DeleteTodoByID(ctx context.Context, id int64) error {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		qry, args := todotodo_storeTemplateImpl.DeleteTodoByID(id)
		slog.Debug("DeleteTodoByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		_, err := db.ExecContext(ctx, qry, args...)
		return err
	}
	return errCtxNoDB
}

// GetTodoByID implements todo_store.Todotodo_store.
func (t *TodoStoreImpl) GetTodoByID(ctx context.Context, id int64) (*todo_store.Todo, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		todo := todo_store.Todo{}
		qry, args := todotodo_storeTemplateImpl.GetTodoByID(id)
		slog.Debug("GetTodoByID", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			err = todotodo_storeTemplateImpl.ErrorMapFind(err)
		}
		return &todo, err
	}
	return nil, errCtxNoDB
}

// FindTodo implements todo_store.Todotodo_store.
func (*TodoStoreImpl) FindTodo(ctx context.Context, filter todo_store.TodoFilter, skip int64, count int) ([]*todo_store.Todo, int64, error) {
	if db, ok := ctx.Value(store.StoreCtxDB).(*sql.DB); ok {
		total := int64(0)
		qry, args := todotodo_storeTemplateImpl.FindTodoTotal(filter)
		slog.Debug("FindTodoTotal", "qry", qry, "args", &store.JsonLogValue{V: args})
		err := db.QueryRowContext(ctx, qry, args...).Scan(&total)
		if err != nil {
			err = todotodo_storeTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		qry, args = todotodo_storeTemplateImpl.FindTodo(filter, skip, count)
		slog.Debug("FindTodo", "qry", qry, "args", &store.JsonLogValue{V: args})
		rows, err := db.QueryContext(ctx, qry, args...)
		if err != nil {
			err = todotodo_storeTemplateImpl.ErrorMapFind(err)
			return nil, total, err
		}
		defer rows.Close()
		todos := []*todo_store.Todo{}
		for rows.Next() {
			todo := todo_store.Todo{}
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
			if err != nil {
				err = todotodo_storeTemplateImpl.ErrorMapFind(err)
				return nil, total, err
			}
			todos = append(todos, &todo)
		}
		return todos, total, nil
	}
	return nil, 0, errCtxNoDB
}
