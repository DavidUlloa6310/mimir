package database

import (
	"testing"

	"github.com/google/uuid"
)

func TestGetAcceleratorByID(t *testing.T) {
	validID := uuid.New().String()
	accelerator, err := GetAcceleratorByID(validID)
	if err != nil {
		t.Errorf("GetAcceleratorByID failed with valid ID: %v", err)
	}
	if accelerator == nil {
		t.Error("GetAcceleratorByID returned nil for valid ID")
	}

	invalidID := "invalid-id"
	accelerator, err = GetAcceleratorByID(invalidID)
	if err == nil {
		t.Error("GetAcceleratorByID should have returned an error for invalid ID")
	}
	if accelerator != nil {
		t.Error("GetAcceleratorByID should have returned nil for invalid ID")
	}
}