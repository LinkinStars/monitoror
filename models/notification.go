package models

type (
	NotificationType string
)

const (
	DingDingNotice NotificationType = "DING_DING"
)

type Notification interface {
	//GetVariantsNames return variant list extract from config
	GetVariantsNames() []VariantName

	//Validate test if config variant is valid
	// return false if empty and error if config have an error (ex: wrong url format)
	Validate(variantName VariantName) (bool, []error)

	//Enable enable notification
	Enable(variantName VariantName)

	// Notify send notification message
	Notify(content string) (err error)

	// Handle handle notification message
	Handle()
}
