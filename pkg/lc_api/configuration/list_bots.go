package configuration

type listBotsRequest struct {
	All bool `json:"all"`
}

type listBotsResponse struct {
	Id string `json:"Id"`
}

func newListBotsRequest() listBotsRequest {
	return listBotsRequest{
		All: true,
	}
}
