package store

import (
	"errors"

	"github.com/go-webapi-team/task-app/entity"
)

var (
	Tags = &TagStore{Tags: map[entity.TagID]*entity.Tag{}}
	ErrNotFound = errors.New("not found")
)

type TagStore struct {
	LastID entity.TagID
	Tags	map[entity.TagID]*entity.Tag
}

func (ts *TagStore) Create(t *entity.Tag) (*entity.TagID, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tags[t.ID] = t
	return &t.ID, nil
}

func (ts *TagStore) Get(id) (*entity.TagID) (*entity.Tag, error) {
	if t, ok := ts.Tags[*id]; ok {
		return t, nil
	}
	return nil, ErrNotFound
}

func (ts *TagStore) GetAll() entity.Tags {
	var tags entity.Tags
	for _, t := range ts.Tags {
		tags = append(tags, t)
	}
	return tags
}