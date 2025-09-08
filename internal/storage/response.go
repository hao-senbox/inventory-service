package storage

type StorageNodeResponse struct {
	Storage  `json:",inline" bson:",inline"`
	Children []*StorageNodeResponse `json:"children,omitempty"`
}

