package app

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"

	"github.com/yokecd/yoke/pkg/flight"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"

	"github.com/tomromeo/yoke-test/pkg/yoke-lib/deployment"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/ingress"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/pv"
	"github.com/tomromeo/yoke-test/pkg/yoke-lib/service"
)

type App struct {
	Labels    map[string]string
	Namespace string
	Name      string

	AddDefaultLabels bool

	//Deployments map[string]*appsv1.Deployment
	//Services    map[string]*corev1.Service
	//Ingresses   map[string]*networkingv1.Ingress
	Resources []flight.Resource
}

func New(name, namespace string, labels map[string]string) *App {
	return &App{
		Labels:           labels,
		Namespace:        namespace,
		Name:             name,
		AddDefaultLabels: true,
		//Deployments:      make(map[string]*appsv1.Deployment),
		//Services:         make(map[string]*corev1.Service),
		//Ingresses:        make(map[string]*networkingv1.Ingress),
		Resources: make([]flight.Resource, 0),
	}
}

func (app *App) PrintResources() []flight.Resource {
	return app.Resources
}

func (app *App) GetNamedResource(name string, gvk schema.GroupVersionKind) (flight.Resource, bool) {
	for _, res := range app.Resources {
		if res.GroupVersionKind() == gvk && res.GetName() == name {
			return res, true
		}
	}
	return nil, false
}

func (app *App) GetResourcesWithKind(gvk schema.GroupVersionKind) []flight.Resource {
	resources := []flight.Resource{}
	for _, res := range app.Resources {
		if res.GroupVersionKind() == gvk {
			resources = append(resources, res)
		}
	}
	return resources
}

func (app *App) WithDeployment(config deployment.Config) *App {
	if config.Name == "" {
		config.Name = app.Name
	}
	if config.Namespace == "" {
		config.Namespace = app.Namespace
	}
	if app.AddDefaultLabels {
		if config.Labels == nil {
			config.Labels = make(map[string]string)
		}
		maps.Copy(config.Labels, app.Labels)
	}

	appDeployment := deployment.New(config)
	app.Resources = append(app.Resources, appDeployment)
	return app
}

func (app *App) WithService(config service.Config) *App {
	if config.Name == "" {
		config.Name = app.Name
	}
	if config.Namespace == "" {
		config.Namespace = app.Namespace
	}
	if app.AddDefaultLabels {
		if config.Labels == nil {
			config.Labels = make(map[string]string)
		}
		maps.Copy(config.Labels, app.Labels)
	}

	appService := service.New(config)
	app.Resources = append(app.Resources, appService)
	return app
}

func (app *App) WithSimpleIngress(config ingress.Config) *App {
	if config.Name == "" {
		config.Name = app.Name
	}
	if config.Namespace == "" {
		config.Namespace = app.Namespace
	}
	if app.AddDefaultLabels {
		if config.Labels == nil {
			config.Labels = make(map[string]string)
		}
		maps.Copy(config.Labels, app.Labels)
	}
	var targetService *corev1.Service
	if config.Service == "" {
		// if there is only one service, take it's name
		services := app.GetResourcesWithKind(schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Service",
		})
		if len(services) == 1 {
			targetService = services[0].(*corev1.Service)
			config.Service = targetService.Name
		} else {
			// if the service name is not specified and there are multiple services, assume the default one uses the app.Name
			if svc, ok := app.GetNamedResource(app.Name, schema.GroupVersionKind{
				Group:   "core",
				Version: "v1",
				Kind:    "Service",
			}); ok {
				targetService = svc.(*corev1.Service)
				config.Service = app.Name
			}
		}
	}
	if targetService == nil {
		panic("Could not find target service for ingress")
	}

	appIngress := ingress.New(config)

	// configure route
	targetPort := targetService.Spec.Ports[0].Port
	appIngress.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: config.HostName,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: ptr.To(networkingv1.PathTypeImplementationSpecific),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: config.Service,
									Port: networkingv1.ServiceBackendPort{
										Number: targetPort,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	app.Resources = append(app.Resources, appIngress)
	return app
}

func (app *App) WithPersistentVolume(config pv.Config) *App {
	if config.Name == "" {
		config.Name = app.Name
	}
	if config.Namespace == "" {
		config.Namespace = app.Namespace
	}
	if app.AddDefaultLabels {
		if config.Labels == nil {
			config.Labels = make(map[string]string)
		}
		maps.Copy(config.Labels, app.Labels)
	}

	pvol, pvolclaim := pv.New(config)

	// bind persistent volume if only one deployment exists
	deployments := app.GetResourcesWithKind(schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	})
	if len(deployments) == 1 {
		deploy := deployments[0].(*appsv1.Deployment)
		deploymentVolumes := deploy.Spec.Template.Spec.Volumes
		deploymentVolumes = append(deploymentVolumes, corev1.Volume{
			Name: "persistence",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvolclaim.Name,
				},
			},
		})
		deploy.Spec.Template.Spec.Volumes = deploymentVolumes
		volumeMounts := deploy.Spec.Template.Spec.Containers[0].VolumeMounts
		if volumeMounts == nil {
			volumeMounts = []corev1.VolumeMount{}
		}
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:              "persistence",
			ReadOnly:          false,
			RecursiveReadOnly: nil,
			MountPath:         config.ContainerPath,
			SubPath:           config.VolumePath,
			MountPropagation:  nil,
			SubPathExpr:       "",
		})
		deploy.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	}

	if config.CreatePV {
		app.Resources = append(app.Resources, pvol)
	}
	app.Resources = append(app.Resources, pvolclaim)
	return app
}

func (app *App) Run() {
	err := json.NewEncoder(os.Stdout).Encode(app.PrintResources())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
