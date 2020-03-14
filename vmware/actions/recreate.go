package actions

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
)

func Recreate(body helpers.Create, uuid uuid.UUID) {
	err := Delete(body.Identifier, body.TargetName, uuid)
	if err != nil {
		return
	}
	Create(body, uuid)
}
