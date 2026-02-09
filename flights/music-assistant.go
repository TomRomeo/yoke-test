package main

import (
	"fmt"
	"io"
	"maps"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	musicAssistantV1 "github.com/tomromeo/yoke-test/crds/music-assistant/v1"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/app"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/deployment"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/ingress"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/pv"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/service"
)

func run() error {

	var musicAssistant musicAssistantV1.MusicAssistant
	if err := yaml.NewYAMLToJSONDecoder(os.Stdin).Decode(&musicAssistant); err != nil && err != io.EOF {
		return err
	}

	// Make sure that our labels include our custom selector.
	if musicAssistant.Spec.Labels == nil {
		musicAssistant.Spec.Labels = map[string]string{}
	}
	labels := musicAssistant.Spec.Labels
	maps.Copy(labels, map[string]string{"app": musicAssistant.Name})
	musicAssistant.Spec.Labels = labels

	application := app.New(musicAssistant.Name, musicAssistant.Namespace, musicAssistant.Spec.Labels).
		WithDeployment(deployment.Config{
			Image:    "ghcr.io/music-assistant/server:2.8.0.dev2026020205",
			Replicas: 1,
		}).
		WithService(service.Config{
			Port:       int32(musicAssistant.Spec.ServicePort),
			TargetPort: int32(musicAssistant.Spec.ServicePort),
		}).
		WithSimpleIngress(ingress.Config{
			IngressClassName: "ingress",
			DefaultBackend:   nil,
			HostName:         musicAssistant.Spec.Hostname,
		}).
		WithPersistentVolume(pv.Config{
			StorageSize: musicAssistant.Spec.StorageSize,
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			StorageClass:  "default",
			NodeAffinity:  nil,
			CreatePV:      true,
			ContainerPath: "/data",
			VolumePath:    "",
		})

	application.Run()

	return nil
}

func main() {

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
