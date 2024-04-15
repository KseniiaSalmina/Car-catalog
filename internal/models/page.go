package models

type CarsPage struct {
	Cars        []Car `json:"cars"`
	PageNo      int   `json:"page_number"`
	Limit       int   `json:"limit"`
	PagesAmount int   `json:"pages_amount"`
}
