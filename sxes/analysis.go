package sxes

import "encoding/json"

type JsonAnalysis struct {
	Uuid            string   `json:"uuid"`
	AcquisitionDate string   `json:"acquisitionDate"`
	Comment         string   `json:"comment"`
	Name            string   `json:"name"`
	Operator        string   `json:"operator"`
	Instrument      string   `json:"instrument"`
	Positions       []string `json:"positionUuids"`
	Images          []string `json:"imageUuids"`
}

type Analysis struct {
	Uuid            string
	AcquisitionDate string
	Comment         string
	Name            string
	Operator        string
	Instrument      string
	Positions       []*Position
	Images          []*Image
}

func MakeAnalysis(uuid string) *Analysis {
	return &Analysis{
		Uuid:      uuid,
		Positions: []*Position{},
		Images:    []*Image{},
	}
}

func (analysis *Analysis) Serialize() (data *JsonAnalysis) {
	data = &JsonAnalysis{}
	data.Uuid = analysis.Uuid
	data.Name = analysis.Name
	data.Comment = analysis.Comment
	data.AcquisitionDate = analysis.AcquisitionDate
	data.Operator = analysis.Operator
	data.Instrument = analysis.Instrument
	data.Positions = []string{}
	data.Images = []string{}

	for _, position := range analysis.Positions {
		data.Positions = append(data.Positions, position.Uuid)
	}

	for _, image := range analysis.Images {
		data.Images = append(data.Images, image.Uuid)
	}

	return
}

func (analysis *Analysis) MarshalJSON() ([]byte, error) {
	return json.Marshal(analysis.Serialize())
}

func (analysis *Analysis) UnmarshalJSON(jsonData []byte) (err error) {
	var data = JsonAnalysis{}
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	analysis.Uuid = data.Uuid
	analysis.Name = data.Name
	analysis.Comment = data.Comment
	analysis.AcquisitionDate = data.AcquisitionDate
	analysis.Operator = data.Operator
	analysis.Instrument = data.Instrument
	analysis.Positions = []*Position{}
	analysis.Images = []*Image{}

	for _, uuid := range data.Positions {
		analysis.Positions = append(analysis.Positions, MakePosition(uuid))
	}

	for _, uuid := range data.Images {
		analysis.Images = append(analysis.Images, MakeImage(uuid))
	}

	return
}
