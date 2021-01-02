package environment

type Zip struct {
	Uri string
}

type Spectrum struct {
	FromZip      string
	FromSpectrum string
	Slice        [2]int
	TypeName     string
	Uuid         string
}

type Environment struct {
	zips     map[string]*Zip
	spectrum map[string]*Spectrum
}

func MakeEnvironment() *Environment {
	return &Environment{
		zips:     map[string]*Zip{},
		spectrum: map[string]*Spectrum{},
	}
}

func (env *Environment) GetZip(name string) *Zip {
	return env.zips[name]
}

func (env *Environment) SetZip(name string, uri string) {
	env.zips[name] = &Zip{
		Uri: uri,
	}
}

func (env *Environment) GetSpectrum(name string) *Spectrum {
	return env.spectrum[name]
}

func (env *Environment) SetSpectrumFromZip(name string, zipName string, typeName string, uuid string) {
	env.spectrum[name] = &Spectrum{
		FromZip:  zipName,
		TypeName: typeName,
		Uuid:     uuid,
	}
}

func (env *Environment) SetSpectrumFromSpectrum(name string, spectrumName string, typeName string, uuid string, slice [2]int) {
	env.spectrum[name] = &Spectrum{
		FromSpectrum: spectrumName,
		Slice:        slice,
		TypeName:     typeName,
		Uuid:         uuid,
	}
}

func (env *Environment) Copy() *Environment {
	zipCopy := map[string]*Zip{}
	for name, z := range env.zips {
		zipCopy[name] = z.Copy()
	}

	specCopy := map[string]*Spectrum{}
	for name, spec := range env.spectrum {
		specCopy[name] = spec.Copy()
	}

	return &Environment{
		zips:     zipCopy,
		spectrum: specCopy,
	}
}

func (z Zip) Copy() *Zip {
	return &Zip{
		Uri: z.Uri,
	}
}

func (spec Spectrum) Copy() *Spectrum {
	return &Spectrum{
		FromZip:      spec.FromZip,
		FromSpectrum: spec.FromSpectrum,
		Slice:        spec.Slice,
		TypeName:     spec.TypeName,
		Uuid:         spec.Uuid,
	}
}
