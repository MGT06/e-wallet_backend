package dto

type FindAllPeople struct {
	ID        int64   `json:"id"`
	Fullname  string  `json:"fullname"`
	ImageURL  *string `json:"image_url"`
	Phone     *string `json:"phone"`
	IsFavorit bool    `json:"is_favorit"`
}

type FindAllPeopleResponse struct {
	TotalData int64           `json:"total_data"`
	SisaData  int64           `json:"sisa_data"`
	Data      []FindAllPeople `json:"data"`
}

type FindAllPeopleRequest struct {
	Fullname string `form:"-"`
	Phone    string `form:"-"`
	Search   string `form:"search"`
	Limit    string `form:"limit"`
}
