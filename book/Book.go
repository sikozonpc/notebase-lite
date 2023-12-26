package book

import "github.com/sikozonpc/notebase/types"

func New(ISBN, title, authors string) *types.Book {
	return &types.Book{
		ISBN:   ISBN,
		Title:  title,
		Authors: authors,
	}
}
