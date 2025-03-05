package main

import (
	"encoding/json"
	"fmt"

	onelog "github.com/francoispqt/onelog"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const httpBadRequestStatusCode = 400

func validate(payload []byte) ([]byte, error) {
	// 从传入的 payload 创建 ValidationRequest 实例
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	// 从 ValidationRequest 对象创建 Settings 实例
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	// 获取原始的 Service JSON 数据
	serviceJSON := validationRequest.Request.Object

	// 尝试将 RAW JSON 解析为 Service 对象
	service := &corev1.Service{}
	if err = json.Unmarshal([]byte(serviceJSON), service); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode Service object: %s", err.Error())),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	logger.DebugWithFields("validating service object", func(e onelog.Entry) {
		e.String("name", service.Metadata.Name)
		e.String("namespace", service.Metadata.Namespace)
	})

	// 检查服务类型是否为 NodePort
	if service.Spec.Type == "NodePort" {
		// 如果设置禁用了 NodePort
		if !settings.IsNodePortAllowed() {
			logger.InfoWithFields("rejecting service object", func(e onelog.Entry) {
				e.String("name", service.Metadata.Name)
				e.String("type", string(service.Spec.Type))
				e.Bool("disable_nodeport", settings.DisableNodePort)
			})

			return kubewarden.RejectRequest(
				kubewarden.Message(
					fmt.Sprintf("Service '%s' of type NodePort is not allowed as per policy configuration", service.Metadata.Name)),
				kubewarden.NoCode)
		}
	}

	return kubewarden.AcceptRequest()
}
