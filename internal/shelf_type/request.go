package shelftype

type CreateShelfTypeRequest struct {
	Name  string  `json:"name" validate:"required"`
	Note  *string `json:"note"`
	Slot  *int    `json:"slot" validate:"required"`
	Level *int    `json:"level" validate:"required"`
}

type UpdateShelfTypeRequest struct {
	Name  string  `json:"name" validate:"required"`
	Note  *string `json:"note"`
	Slot  *int    `json:"slot" validate:"required"`
	Level *int    `json:"level" validate:"required"`
}
