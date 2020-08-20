package update

import (
	"testing"

	argov1alpha1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

const (
	testImage1 = "docker.io/bigkevmcd/go-demo:af93dae"
	testImage2 = "docker.io/bigkevmcd/go-demo:72ab9cc"
)

func TestOverrideImages(t *testing.T) {
	updateTests := []struct {
		desc    string
		initial argov1alpha1.KustomizeImages
		want    argov1alpha1.KustomizeImages
	}{
		{
			desc:    "empty set",
			initial: []argov1alpha1.KustomizeImage{},
			want:    []argov1alpha1.KustomizeImage{testImage1},
		},
		{
			desc:    "existing image - different tag",
			initial: []argov1alpha1.KustomizeImage{testImage2},
			want:    []argov1alpha1.KustomizeImage{testImage1},
		},
		{
			desc:    "existing image - same tag",
			initial: []argov1alpha1.KustomizeImage{testImage1},
			want:    []argov1alpha1.KustomizeImage{testImage1},
		},
	}

	for _, tt := range updateTests {
		kust := &argov1alpha1.ApplicationSourceKustomize{
			Images: tt.initial,
		}
		app := argov1alpha1.Application{
			Spec: argov1alpha1.ApplicationSpec{
				Source: argov1alpha1.ApplicationSource{
					Kustomize: kust,
				},
			},
		}

		OverrideImage(&app, testImage1)

		if diff := cmp.Diff(tt.want, kust.Images); diff != "" {
			t.Errorf("%s failed comparison:\n%s", tt.desc, diff)
		}
	}
}
