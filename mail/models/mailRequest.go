package models

type MailRequest struct {
	// To        string   `form:"to"`
	To        []string `form:"to" binding:"required"`
	FromEmail string   `form:"fromEmail"`
	Subject   string   `form:"subject"`
	Body      string   `form:"body"`
	Body1     string   `form:"body1"`
	Body2     string   `form:"body2"`
	BodyLink  string   `form:"bodylink"`
	LinkName  string   `form:"linkname"`
	CC        []string `form:"cc"`
}
