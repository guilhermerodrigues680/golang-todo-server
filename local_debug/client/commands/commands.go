package commands

type Command string

const (
	Readall        Command = "readall"
	Create         Command = "create"
	Read           Command = "read"
	Update         Command = "update"
	Delete         Command = "delete"
	CreateMultiple Command = "create-multiple"
	DeleteMultiple Command = "delete-multiple"
)

func (c Command) Cmd() string {
	return string(c)
}
