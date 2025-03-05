package main

import (
	"encoding/json"
	"testing"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := json.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if settings.DisableNodePort {
		t.Errorf("Expected DisableNodePort to be false by default")
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

func TestNodePortControl(t *testing.T) {
	// 测试场景1：禁用 NodePort
	settings := Settings{
		DisableNodePort: true,
	}

	if settings.IsNodePortAllowed() {
		t.Errorf("NodePort should be disabled when DisableNodePort is true")
	}

	// 测试场景2：允许 NodePort
	settings = Settings{
		DisableNodePort: false,
	}

	if !settings.IsNodePortAllowed() {
		t.Errorf("NodePort should be allowed when DisableNodePort is false")
	}
}

func TestParsingSettingsWithValueProvided(t *testing.T) {
	rawSettings := []byte(`{"disable_nodeport": true}`)
	settings := &Settings{}
	if err := json.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if !settings.DisableNodePort {
		t.Errorf("Expected DisableNodePort to be true")
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}
