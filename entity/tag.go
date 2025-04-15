package entity

import "time"

type TagID int64
type UserID int64

type Tag struct {
	ID      TagID      `json:"id"`
	TagName string     `json:"name"`
	UserID  UserID 	   `json:"user_id" `
	Created time.Time  `json:"created"`
}

type Tags []*Tag