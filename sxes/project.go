package sxes

import "encoding/json"

type JsonProject struct {
	Uuid     string   `json:"uuid"`
	Name     string   `json:"name"`
	Comment  string   `json:"comment"`
	Analyses []string `json:"analyses"`
}

type Project struct {
	Uuid     string
	Name     string
	Comment  string
	Analyses []*Analysis
}

func MakeProject(uuid string) *Project {
	return &Project{
		Uuid:     uuid,
		Analyses: []*Analysis{},
	}
}

func (project *Project) Serialize() (data *JsonProject) {
	data = &JsonProject{}
	data.Uuid = project.Uuid
	data.Name = project.Name
	data.Comment = project.Comment

	for _, analysis := range project.Analyses {
		data.Analyses = append(data.Analyses, analysis.Uuid)
	}

	return
}

func (project *Project) MarshalJSON() ([]byte, error) {
	return json.Marshal(project.Serialize())
}

func (project *Project) UnmarshalJSON(jsonData []byte) (err error) {
	var data = JsonProject{}
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	for _, uuid := range data.Analyses {
		project.Analyses = append(project.Analyses, MakeAnalysis(uuid))
	}

	project.Uuid = data.Uuid
	project.Name = data.Name
	project.Comment = data.Comment

	return
}
