package facades

import (
	"github.com/eddyjj92/goravel-inertia/contracts"
)

var instance contracts.Inertia

func Inertia() contracts.Inertia {
	return instance
}

func RegisterInertia(inertia contracts.Inertia) {
	instance = inertia
}
