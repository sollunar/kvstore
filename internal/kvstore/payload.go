package kvstore

type SetRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
