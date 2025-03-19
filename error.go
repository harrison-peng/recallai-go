package recallaigo

type Error struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func (e Error) Error() string {
	return e.Detail
}
