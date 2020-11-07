package main

import "strconv"

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

func (t *token) tokenToType() string {
	if t.Type == 0 {
		return "PING"
	}
	return "PONG"
}

func (t *token) tokenToValue() string {
	if t != nil {
		return strconv.Itoa(t.Value)
	}
	return "_"
}
