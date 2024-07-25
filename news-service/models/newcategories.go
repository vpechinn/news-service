package models

type NewsCategory struct {
	NewsID     int64 `reform:"NewsId"`
	CategoryID int64 `reform:"CategoryId"`
}
