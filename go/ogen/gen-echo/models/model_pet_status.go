package models
// PetStatus : pet status in the store
type PetStatus string

// List of PetStatus
const (
	AVAILABLE PetStatus = "available"
	PENDING PetStatus = "pending"
	SOLD PetStatus = "sold"
)
