package model

import "time"

type RequestCode struct {
	Login           string `json:"login"`
	Scope           string `json:"scope"`
	TillNextRequest int64  `json:"till_next_request,omitempty"`
}

type TokenInfo struct {
	DocId       string
	Id          string    `firestore:"id"`
	VkId        int64     `firestore:"vkid"`
	Scope       string    `firestore:"scope"`
	Token       string    `firestore:"token"`
	CreateToken time.Time `firestore:"create_token"`
}

type CodeInfo struct {
	Code       int64     `firestore:"code"`
	SendCodeAt time.Time `firestore:"send_code_at"`
	TryCount   int64     `firestore:"try_count"`
}

type CheckCodeRequest struct {
	Login string `json:"login"`
	Code  int64  `json:"code"`
	Scope string `json:"scope"`
}
