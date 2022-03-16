package representations

type CreateBlockchainInput struct {
	To string `json:"to" binding:"required"`
}