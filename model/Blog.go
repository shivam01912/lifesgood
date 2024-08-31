package model

type Blog struct {
	Title        string
	Brief        string
	Tags         []string
	Content      []byte
	Views        int `default:0`
	Likes        int `default:0`
	Comments     []string
	CreatedAt    int64
	ModifiedDate int64
}
