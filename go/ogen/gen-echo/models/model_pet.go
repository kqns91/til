package models

type Pet struct {

	Id int64 `json:"id,omitempty"`

	Name string `json:"name"`

	PhotoUrls []string `json:"photoUrls,omitempty"`

	Status PetStatus `json:"status,omitempty"`
}
