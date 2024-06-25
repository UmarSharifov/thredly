package main

import (
	"html/template"

	"thredly.com/thredly/pkg/models"
)

type templateData struct {
	User                   *models.User
	Users                  []*models.User
	Tred                   *models.Tred
	Treds                  []*models.Tred
	TredsWithChilds        []*models.TredsWithChilds
	Categories             []*models.Categories
	Events                 []*models.Events
	EventsCategories       []*models.EventCategories
	CurrentCategory        string
	CurrentMenu            template.HTML
	ActiveValue            string
	ComplaintWithDetails   []*models.ComplaintWithDetails
	ErrorMessageCreateTred string
	CurrentTredBackstage   string
	UserSubsTo             int
	UserSubsFrom           int
	UserTredsCount         int
	ParentTredId           string
}
