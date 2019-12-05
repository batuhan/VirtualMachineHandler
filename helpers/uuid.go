package helpers

import (
	"github.com/google/uuid"
	"log"
)

func GenerateUUID() uuid.UUID {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err.Error())
	}
	return newUUID
}
