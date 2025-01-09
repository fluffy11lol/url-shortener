package storage

import (
	"errors"
)

var (
	ErrUrlNotFound     = errors.New("url not found")
	ErrUrlAlreadyExist = errors.New("url already exist")
)

type Storage interface {
	SaveAlias(alias, url string) error
	GetURL(alias string) (string, error)
}
