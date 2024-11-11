package entities

type Observer interface {
	Update(status string)
	GetID() int
}
