package response

type TeamShowItem struct {
	ID          int    `json:"id"`
	AvatarUrl   string `json:"avatarUrl"`
	Name        string `json:"name"`
	IsActivated int    `json:"isActivated" `
	IsCaptain   int    `json:"isCaptain" `
}
type TeamShowResponse struct {
	JoinedTeam []TeamShowItem   `json:"joinedTeam"`
	MyTeams    [][]TeamShowItem `json:"myTeams"`
}

type ConsumeData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

//	type ConsumeDetails struct {
//		Name string        `json:"name"`
//		Data []ConsumeData `json:"data"`
//	}
type ConsumeDetailsResponse struct {
	Details map[string][]ConsumeData `json:"details"`
}
