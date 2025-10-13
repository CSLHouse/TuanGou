package request

type DetailsInfo struct {
	UserIds []int `json:"userIds" binding:"required"`
}
