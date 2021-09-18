package email

type Template string

const (
	AccountCreated Template = "account_created"
	ItemSelected            = "item_selected"
	ItemApproved            = "item_approved"
	Winner                  = "winner"
	Losener                 = "looser"
)
