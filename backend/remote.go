package backend

import (
	"PLViewer/sxes"
	"encoding/json"
	"errors"
	"github.com/bob620/baka-rpc-go/parameters"
)

func (backend Backend) Load(uri string) error {
	_, resErr := backend.Client.CallMethodByName(nil, "load", &parameters.StringParam{Name: "uri", Default: uri})
	if resErr != nil {
		return errors.New(resErr.Message)
	}
	return nil
}

func (backend Backend) GetProjects() (projects []*sxes.Project, err error) {
	rawData, resErr := backend.Client.CallMethodWithNone(nil, "getProjects")
	if resErr != nil {
		return nil, errors.New(resErr.Message)
	}
	test, _ := rawData.MarshalJSON()
	projects = []*sxes.Project{}
	err = json.Unmarshal(test, &projects)
	if err != nil {
		return nil, err
	}
	return projects, err
}

func (backend Backend) GetAnalysis(uuid string) (analysis *sxes.Analysis, err error) {
	rawData, resErr := backend.Client.CallMethodByPosition(nil, "getAnalysis", &parameters.StringParam{Name: "uuid", Default: uuid})
	if resErr != nil {
		return nil, errors.New(resErr.Message)
	}
	test, _ := rawData.MarshalJSON()
	analysis = &sxes.Analysis{}
	err = json.Unmarshal(test, &analysis)
	if err != nil {
		return nil, err
	}
	return analysis, err
}

func (backend Backend) GetPosition(uuid string) (position *sxes.Position, err error) {
	rawData, resErr := backend.Client.CallMethodByPosition(nil, "getPosition", &parameters.StringParam{Name: "uuid", Default: uuid})
	if resErr != nil {
		return nil, errors.New(resErr.Message)
	}
	test, _ := rawData.MarshalJSON()
	position = &sxes.Position{}
	err = json.Unmarshal(test, &position)
	if err != nil {
		return nil, err
	}
	return position, err
}

func (backend Backend) GetLine(uuid string, typeName string) (line []float64, err error) {
	rawData, resErr := backend.Client.CallMethodByPosition(nil, "getLine", &parameters.StringParam{Name: "uuid", Default: uuid}, &parameters.StringParam{Name: "typeName", Default: typeName})
	if resErr != nil {
		return nil, errors.New(resErr.Message)
	}
	test, _ := rawData.MarshalJSON()
	lines := [][]float64{}
	err = json.Unmarshal(test, &lines)
	if err != nil {
		return nil, err
	}
	return lines[0], err
}
