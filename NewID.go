package main

import "github.com/google/uuid"

func NewID() string {
	newId, _ := uuid.NewUUID()
	return newId.String()
}
