package database

import (
	"fmt"
	"testing"
)

func TestGetAllAcceleratorsAndIndividually(t *testing.T) {
    acceleratorIDs, err := GetAllAccelerators()
    if err != nil {
        t.Fatalf("Failed to get all accelerators: %v", err)
    }

    if len(acceleratorIDs) == 0 {
        t.Log("No accelerators found in the database. Skipping individual tests.")
        return
    }

    for _, id := range acceleratorIDs {
        t.Run(fmt.Sprintf("TestGetAcceleratorByID_%s", id), func(t *testing.T) {
            accelerator, err := GetAcceleratorByID(id)
            if err != nil {
                t.Errorf("Failed to get accelerator with ID %s: %v", id, err)
                return
            }

            if accelerator == nil {
                t.Errorf("GetAcceleratorByID returned nil for ID %s", id)
                return
            }
        })
    }
}