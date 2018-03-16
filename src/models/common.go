package models

type MapInterface map[string]interface{}
type MapString map[string]string

type JsonResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func NewJsonResponse(data, meta interface{}) *JsonResponse {
	return &JsonResponse{
		Data: data,
		Meta: meta,
	}
}
