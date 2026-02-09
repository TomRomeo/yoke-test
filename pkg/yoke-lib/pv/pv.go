package pv

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type Config struct {
	Name                   string
	Namespace              string
	Labels                 map[string]string
	StorageSize            string
	AccessModes            []corev1.PersistentVolumeAccessMode
	StorageClass           string
	NodeAffinity           *corev1.VolumeNodeAffinity
	CreatePV               bool
	ContainerPath          string
	VolumePath             string
	PersistentVolumeSource corev1.PersistentVolumeSource
}

func New(cfg Config) (*corev1.PersistentVolume, *corev1.PersistentVolumeClaim) {
	var pv *corev1.PersistentVolume
	if cfg.CreatePV {
		pv = &corev1.PersistentVolume{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PersistentVolume",
				APIVersion: corev1.SchemeGroupVersion.Identifier(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cfg.Name,
				Namespace: cfg.Namespace,
			},
			Spec: corev1.PersistentVolumeSpec{
				Capacity: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cfg.StorageSize),
				},
				PersistentVolumeSource: corev1.PersistentVolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/tmp/pv"},
				},
				AccessModes:                   cfg.AccessModes,
				ClaimRef:                      nil,
				PersistentVolumeReclaimPolicy: "",
				StorageClassName:              cfg.StorageClass,
				MountOptions:                  nil,
				VolumeMode:                    nil,
				NodeAffinity:                  cfg.NodeAffinity,
				VolumeAttributesClassName:     nil,
			},
		}
	}

	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: cfg.AccessModes,
			Selector:    nil,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cfg.StorageSize),
				},
			},
			VolumeName:                cfg.Name,
			StorageClassName:          ptr.To(cfg.StorageClass),
			VolumeMode:                nil,
			DataSource:                nil,
			DataSourceRef:             nil,
			VolumeAttributesClassName: nil,
		},
	}
	return pv, pvc
}
