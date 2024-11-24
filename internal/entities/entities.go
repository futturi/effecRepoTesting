package entities

type Song struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	GroupName   string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongUpdate struct {
	GroupName   *string `json:"group"`
	Song        *string `json:"song"`
	ReleaseDate *string `json:"releaseDate"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
}

type SongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type TextResponse struct {
	Text string `json:"text"`
}

type InsertResponse struct {
	ID int `json:"id"`
}

type UpdateResponse struct {
	ID string `json:"id"`
}

type DeleteResponse struct {
	Status bool `json:"status"`
}
