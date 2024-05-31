package models

type Session struct {
	Token string `json:"token"`
	TTL   int64  `json:"ttl"` // TTL for session expiration
}
