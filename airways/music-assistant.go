package main

import (
	"encoding/json"
	"fmt"
	"os"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/yokecd/yoke/pkg/apis/v1alpha1"
	"github.com/yokecd/yoke/pkg/openapi"

	v1 "github.com/tomromeo/yoke-test/crds/music-assistant/v1"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	return json.NewEncoder(os.Stdout).Encode(v1alpha1.Airway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "music-assistants.torodo.io",
		},
		Spec: v1alpha1.AirwaySpec{
			Mode: v1alpha1.AirwayModeStandard,
			WasmURLs: v1alpha1.WasmURLs{
				Flight: "https://github.com/tomromeo/yoke-test/releases/download/latest/music-assistant_v1_flight.wasm.gz",
			},
			Template: apiextv1.CustomResourceDefinitionSpec{
				Group: "torodo.io",
				Names: apiextv1.CustomResourceDefinitionNames{
					Plural:     "music-assistants",
					Singular:   "music-assistant",
					ShortNames: []string{"ma"},
					Kind:       "MusicAssistant",
				},
				Scope: apiextv1.NamespaceScoped,
				Versions: []apiextv1.CustomResourceDefinitionVersion{
					{
						Name:    "v1",
						Served:  true,
						Storage: true,
						Schema: &apiextv1.CustomResourceValidation{
							OpenAPIV3Schema: openapi.SchemaFor[v1.MusicAssistant](),
						},
					},
				},
			},
		},
	})
}
