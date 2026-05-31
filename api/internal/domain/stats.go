package domain

// StudyStats is the response shape for the GET /stats endpoint.
type StudyStats struct {
	Streak StreakStats `json:"streak"`
	Today  DayStats   `json:"today"`
	Total  TotalStats `json:"total"`
}

type StreakStats struct {
	Current  int    `json:"current"`
	Longest  int    `json:"longest"`
	LastDate string `json:"lastDate"` // "2006-01-02" or ""
}

type DayStats struct {
	Reviews int `json:"reviews"`
}

type TotalStats struct {
	Reviews      int `json:"reviews"`
	CardsStudied int `json:"cardsStudied"`
}
