package v1

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	APIVersion         = "torodo.io/v1"
	KindMusicAssistant = "MusicAssistant"
)

type MusicAssistant struct {
	metav1.TypeMeta
	metav1.ObjectMeta `json:"metadata"`
	Spec              MusicAssistantSpec `json:"spec"`
}

// Our MusicAssistant Specification
type MusicAssistantSpec struct {
	Image       string            `json:"image"`
	Replicas    int32             `json:"replicas"`
	Labels      map[string]string `json:"labels,omitempty"`
	NodePort    int               `json:"nodePort,omitempty"`
	ServicePort int               `json:"port,omitempty"`
	Hostname    string            `json:"hostname,omitempty"`
	StorageSize string            `json:"storageSize,omitempty"`
}

// Custom Marshalling Logic so that users do not need to explicity fill out the Kind and ApiVersion.
func (musicAssistant MusicAssistant) MarshalJSON() ([]byte, error) {
	musicAssistant.Kind = KindMusicAssistant
	musicAssistant.APIVersion = APIVersion

	type MusicAssistantAlt MusicAssistant
	return json.Marshal(MusicAssistantAlt(musicAssistant))
}

// Custom Unmarshalling to raise an error if the ApiVersion or Kind does not match.
func (musicAssistant *MusicAssistant) UnmarshalJSON(data []byte) error {
	type BackendAlt MusicAssistant
	if err := json.Unmarshal(data, (*BackendAlt)(musicAssistant)); err != nil {
		return err
	}
	if musicAssistant.APIVersion != APIVersion {
		return fmt.Errorf("unexpected api version: expected %s but got %s", APIVersion, musicAssistant.APIVersion)
	}
	if musicAssistant.Kind != KindMusicAssistant {
		return fmt.Errorf("unexpected kind: expected %s but got %s", KindMusicAssistant, musicAssistant.Kind)
	}
	return nil
}
