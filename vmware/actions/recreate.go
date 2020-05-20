package actions

import (
	"VirtualMachineHandler/helpers"
	"github.com/google/uuid"
)

func Recreate(body helpers.Create, uuid uuid.UUID) {
	err := Delete(body.LocationId, body.TargetName, uuid)
	if err != nil {
		return
	}
	Create(body, uuid)
}
