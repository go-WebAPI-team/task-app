package entity

import "time"

type TagID int64
type TagName string
type UserID int64

type Tag struct {
	ID      TaskID     `json:"id"`
	Name   	string     `json:"name"`
	UserID  UserID 	   `json:"userid" `
	Created time.Time  `json:"created"`
}

type Tags []*Tag