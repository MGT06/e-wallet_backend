package dto

type FindAllPeopleResponse struct {
	ID        int64   `json:"id"`
	Fullname  string  `json:"fullname"`
	ImageURL  *string `json:"image_url"`
	Phone     *string `json:"phone"`
	IsFavorit bool    `json:"is_favorit"`
}

type FindAllPeopleRequest struct {
	Fullname string `form:"-"`
	Phone    string `form:"-"`
	Search   string `form:"search"`
}
