package database

import (
	"testing"

	"github.com/davidulloa/mimir/models"
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

func TestCreateAccelerator(t *testing.T) {
	validAccelerator := &models.Accelerator{
		Url:         "http://example.com",
		Title:       "Test Accelerator",
		Description: "This is a test accelerator",
		Category:    "Test",
	}
	_, err := CreateAccelerator(validAccelerator)
	if err != nil {
		t.Errorf("CreateAccelerator failed with valid data: %v", err)
	}

	invalidAccelerator := &models.Accelerator{}
	_, err = CreateAccelerator(invalidAccelerator)
	if err == nil {
		t.Error("CreateAccelerator should have returned an error for invalid data")
	}
}

func TestDeleteAccelerator(t *testing.T) {
	// Test with a valid ID
	validID := uuid.New().String()
	err := DeleteAccelerator(validID)
	if err != nil {
		t.Errorf("DeleteAccelerator failed with valid ID: %v", err)
	}

	// Test with an invalid ID
	invalidID := "invalid-id"
	err = DeleteAccelerator(invalidID)
	if err == nil {
		t.Error("DeleteAccelerator should have returned an error for invalid ID")
	}
}

func TestIntegration(t *testing.T) {
	// Create an accelerator
	newAccelerator := &models.Accelerator{
		Url:         "http://example.com",
		Title:       "Integration Test Accelerator",
		Description: "This is an integration test accelerator",
		Category:    "Test",
	}
	accelerator_id, err := CreateAccelerator(newAccelerator)
	if err != nil {
		t.Fatalf("Failed to create accelerator: %v", err)
	}

	createdAccelerator, err := GetAcceleratorByID(accelerator_id)
	if err != nil {
		t.Fatalf("Failed to get created accelerator: %v", err)
	}
	if createdAccelerator.Title != newAccelerator.Title {
		t.Errorf("Retrieved accelerator does not match created accelerator")
	}

	err = DeleteAccelerator(accelerator_id)
	if err != nil {
		t.Fatalf("Failed to delete accelerator: %v", err)
	}

	// Verify deletion
	deletedAccelerator, err := GetAcceleratorByID(accelerator_id)
	if err == nil || deletedAccelerator != nil {
		t.Error("Accelerator was not deleted successfully")
	}
}