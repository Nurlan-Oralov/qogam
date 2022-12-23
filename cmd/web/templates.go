package main

import "golangify.com/snippetbox/pkg/models"

// Добавляем поле Snippets в структуру templateData
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
