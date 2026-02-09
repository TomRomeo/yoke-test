package main

import (
	"github.com/yokecd/yoke/pkg/flight"
	v1 "k8s.io/api/core/v1"

	"github.com/tomromeo/yoke-test/pkg/yoke-lib/app"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/deployment"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/ingress"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/pv"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/service"
)

func main() {
	defaultName := flight.Release()
	defaultNamespace := flight.Namespace()
	defaultLabels := map[string]string{"app": flight.Release()}

	appPort := int32(8095)
	hostName := "music-assistant.k8s.torodo.io"
	storageSize := "2Mi"

	application := app.New(defaultName, defaultNamespace, defaultLabels).
		WithDeployment(deployment.Config{
			Image:    "ghcr.io/music-assistant/server:2.8.0.dev2026020205",
			Replicas: 1,
		}).
		WithService(service.Config{
			Port:       appPort,
			TargetPort: appPort,
		}).
		WithSimpleIngress(ingress.Config{
			IngressClassName: "ingress",
			DefaultBackend:   nil,
			HostName:         hostName,
		}).
		WithPersistentVolume(pv.Config{
			StorageSize: storageSize,
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
}
