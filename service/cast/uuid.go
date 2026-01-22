package cast

import (
	"log"

	"github.com/google/uuid"
)

func CastUUID(input string) (uuid.UUID, error) {
	newString, err := uuid.Parse(input)
	if err != nil {
		log.Println("invalue uuid: ", err)

		return newString, err
	}

	return newString, nil
}
