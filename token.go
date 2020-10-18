package main

type tokenType int

const (
	ping tokenType = iota
	pong
)

type token struct {
	Type  tokenType `json:"type"`
	Value int       `json:"value"`
}

func newToken() *token {
	return &token{}
}
