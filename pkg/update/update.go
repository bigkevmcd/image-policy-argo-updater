package update

import (
	argov1alpha1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
)

func OverrideImage(a *argov1alpha1.Application, newImage argov1alpha1.KustomizeImage) error {
	images := a.Spec.Source.Kustomize.Images
	images = removeImage(images, newImage)
	a.Spec.Source.Kustomize.Images = append(images, newImage)
	return nil
}

func removeImage(imgs []argov1alpha1.KustomizeImage, img argov1alpha1.KustomizeImage) []argov1alpha1.KustomizeImage {
	updated := []argov1alpha1.KustomizeImage{}
	for _, v := range imgs {
		if !v.Match(img) {
			updated = append(updated, v)
		}
	}
	return updated
}
