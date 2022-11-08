package model

type User struct {
	DocId     string
	FirstName string   `firestore:"firstname"`
	LastName  string   `firestore:"lastname"`
	TgId      string   `firestore:"tgid"`
	VkId      int64    `firestore:"vkid"`
	VkDomain  string   `firestore:"vkdomain"`
	Roles     []string `firestore:"roles"`
}
