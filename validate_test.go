package main

import (
	"encoding/json"
	"testing"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

func TestEmptySettingsLeadsToApproval(t *testing.T) {
	settings := Settings{}
	service := corev1.Service{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: &corev1.ServiceSpec{
			Type: "ClusterIP",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&service, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err = json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, *response.Code)
	}
}

func TestNodePortAllowed(t *testing.T) {
	settings := Settings{
		DisableNodePort: false,
	}
	service := corev1.Service{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: &corev1.ServiceSpec{
			Type: "NodePort",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&service, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err = json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Error("Unexpected rejection")
	}
}

func TestNodePortDisabled(t *testing.T) {
	settings := Settings{
		DisableNodePort: true,
	}
	service := corev1.Service{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: &corev1.ServiceSpec{
			Type: "NodePort",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&service, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err = json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted {
		t.Error("Unexpected approval")
	}

	expectedMessage := "Service 'test-service' of type NodePort is not allowed as per policy configuration"
	if response.Message == nil {
		t.Errorf("expected response to have a message")
	}
	if *response.Message != expectedMessage {
		t.Errorf("Got '%s' instead of '%s'", *response.Message, expectedMessage)
	}
}

func TestNonNodePortService(t *testing.T) {
	settings := Settings{
		DisableNodePort: true,
	}
	service := corev1.Service{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: &corev1.ServiceSpec{
			Type: "ClusterIP",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&service, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err = json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Error("Unexpected rejection")
	}
}
