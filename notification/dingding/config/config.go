package config

type (
	DingDing struct {
		Webhook string `validate:"required,url,http"`
		Secret  string
	}
)

var Default = &DingDing{
	Webhook: "",
	Secret:  "",
}
