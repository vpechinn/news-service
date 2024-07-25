package models

type News struct {
	ID         int64   `reform:"Id,pk"`
	Title      string  `reform:"Title"`
	Content    string  `reform:"Content"`
	Categories []int64 `reform:"-"`
}
