package bots

type CreateBotDto struct {
	Name          string `json:"name"`
	OwnerClientId string `json:"owner_client_id"`
}

func NewDefaultBot(name, ownerClientId string) CreateBotDto {
	return CreateBotDto{
		Name:          name,
		OwnerClientId: ownerClientId,
	}
}
