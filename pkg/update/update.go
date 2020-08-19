package update

import (
	argov1alpha1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
)

// Image contains an image name and a new tag.
type Image struct {
	// Name is a tag-less image name.
	Name string `json:"name,omitempty"`

	// NewTag is the value used to replace the original tag.
	NewTag string `json:"newTag,omitempty"`
}

func OverrideImages(a *argov1alpha1.Application, updates []Image) error {
	return nil
}
