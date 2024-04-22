package configuration

const (
	AGENT_STATUS_ACCEPTING_CHATS = "accepting_chats"
	AGENT_STATUS_NOT_ACCEPTING_CHATS = "not_accepting_chats"
)

type CreateBotResponse struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

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
