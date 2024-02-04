package store

import "github.com/senomas/gosvc_store/store"

type Todo struct {
	Title     string
	ID        int64
	Completed bool
}

type TodoFilter struct {
	Title     store.FilterString
	ID        store.FilterInt64
	Completed store.FilterBool
}
