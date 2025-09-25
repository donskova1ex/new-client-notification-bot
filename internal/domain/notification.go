package domain

type Notification struct {
	Phone            string `json:"phone"`
	CompanyName      string `json:"company_name"`
	NotificationText string `json:"notification_text"`
}
