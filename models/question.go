package models

type Question struct {
	Id       int
	Number   string
	Wording  string
	Addition string
	Verify   []string
}
