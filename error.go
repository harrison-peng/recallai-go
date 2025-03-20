package recallaigo

type Error struct {
	Code   string    `json:"code"`
	Detail string `json:"detail"`
}

func (e Error) Error() string {
	return e.Detail
}
