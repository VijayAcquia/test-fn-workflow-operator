package controller

import (
	"github.com/acquia/fn-workflows-operator/pkg/controller/customerimagebuild"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, customerimagebuild.Add)
}
