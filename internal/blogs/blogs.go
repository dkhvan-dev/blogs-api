package blogs

import (
	"blogs-api/internal/users/model"
)

type Blog struct {
	Id      int        `json:"id"`
	Title   string     `json:"title"`
	Content string     `json:"content"`
	Author  model.User `json:"author"`
}
