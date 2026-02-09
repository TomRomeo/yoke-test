package main

import (
	"github.com/tomromeo/yoke-test/pkg/util"

	"github.com/tomromeo/yoke-test/pkg/yoke-lib/app"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/deployment"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/ingress"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/service"
)

func main() {
	fc := util.NewFlightConfig()
	defaultName := fc.Release
	defaultNamespace := fc.Namespace
	defaultLabels := fc.Labels

	appPort := int32(8095)
	hostName := "music-assistant.k8s.torodo.io"

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
		})

	fc.Resources = application.Resources()
	fc.Run()
}
