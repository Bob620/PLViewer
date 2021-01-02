package sxes

type Position struct {
	Uuid         string   `json:"uuid"`
	Comment      string   `json:"comment"`
	Operator     string   `json:"operator"`
	Background   string   `json:"background"`
	Condition    string   `json:"condition"`
	RawCondition string   `json:"rawCondition"`
	Types        []string `json:"types"`
}

func MakePosition(uuid string) *Position {
	return &Position{
		Uuid: uuid,
	}
}
