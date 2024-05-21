package main

import "thredly.com/thredly/pkg/models"

// Создаем тип templateData, который будет действовать как хранилище для
// любых динамических данных, которые нужно передать HTML-шаблонам.
// На данный момент он содержит только одно поле, но мы добавим в него другие
// по мере развития нашего приложения.
type templateData struct {
	Snippet    *models.Snippet
	Snippets   []*models.Snippet
	User       *models.User
	Users      []*models.User
	Tred       *models.Tred
	Treds      []*models.Tred
	TredsChild []*models.Tred
}
