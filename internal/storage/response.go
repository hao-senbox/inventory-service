package storage

type StorageNodeResponse struct {
	Storage      Storage                `json:"storage"`
	ImageMainUrl string                 `json:"image_main_url"`
	ImageMapUrl  string                 `json:"image_map_url"`
	Children     []*StorageNodeResponse `json:"children"`
}
