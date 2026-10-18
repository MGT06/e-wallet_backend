package dto

type FindAllPeople struct {
	ID        int64   `json:"id"`
	Fullname  string  `json:"fullname"`
	ImageURL  *string `json:"image_url"`
	Phone     *string `json:"phone"`
	IsFavorit bool    `json:"is_favorit"`
}

type FindAllPeopleResponse struct {
	Data []FindAllPeople `json:"data"`
}

type FindAllPeopleRequest struct {
	Fullname string `form:"-"`
	Phone    string `form:"-"`
	Search   string `form:"search"`
	Limit    int    `form:"limit"`
}

type HistoryTransactionRequest struct {
	Fullname string `form:"-"`
	Phone    string `form:"-"`
	Search   string `form:"search"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"page"`
	Type     string `form:"page"`
}

type HistoryTransaction struct {
	ID       int64   `json:"id"`
	Type     string  `json:"type"`
	Amount   int64   `json:"amount"`
	Fullname string  `json:"fullname"`
	ImageURL *string `json:"image_url"`
	Phone    *string `json:"phone"`
}

type HistoryTransactionResponse struct {
	Data []HistoryTransaction `json:"data"`
}
