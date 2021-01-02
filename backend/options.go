package backend

type Options interface {
	Serialize() []string
}

type CreatorOptions struct {
	ExportTo    string
	ConvertXes  bool
	ConvertQlw  bool
	ConvertMap  bool
	ConvertLine bool
	ConvertQMap bool
	Loose       bool
	Recover     bool
	Debug       bool
}

func MakeCreatorOptions() *CreatorOptions {
	return &CreatorOptions{
		ExportTo:    "zip",
		ConvertXes:  true,
		ConvertQlw:  true,
		ConvertMap:  false,
		ConvertLine: false,
		ConvertQMap: false,
		Loose:       false,
		Recover:     false,
		Debug:       false,
	}
}

func (options CreatorOptions) Serialize() []string {
	args := []string{}
	switch options.ExportTo {
	case "csv":
		args = append(args, "-f c")
		break
	case "jeol":
		args = append(args, "-f c", "-j")
		break
	case "json":
		args = append(args, "-f j")
		break
	case "plzip":
		args = append(args, "-f z")
		break
	}

	if options.ConvertXes {
		args = append(args, "-x")
	}
	if options.ConvertQlw {
		args = append(args, "-q")
	}
	if options.ConvertMap {
		args = append(args, "-m")
	}
	if options.ConvertLine {
		args = append(args, "-l")
	}
	if options.ConvertQMap {
		args = append(args, "-k")
	}
	if options.Loose {
		args = append(args, "-y")
	}
	if options.Recover {
		args = append(args, "-r")
	}
	if options.Debug {
		args = append(args, "-d")
	}

	return args
}

type ProcessorOptions struct {
	ExportToZip bool
	ConvertJeol bool
	ConvertXes  bool
	ConvertQlw  bool
	ConvertMap  bool
	ConvertLine bool
	ConvertQMap bool
	Loose       bool
	Recover     bool
	Debug       bool
}

func MakeProcessorOptions() *ProcessorOptions {
	return &ProcessorOptions{
		ExportToZip: true,
		ConvertJeol: false,
		ConvertXes:  true,
		ConvertQlw:  true,
		ConvertMap:  false,
		ConvertLine: false,
		ConvertQMap: false,
		Loose:       false,
		Recover:     false,
		Debug:       false,
	}
}

func (options ProcessorOptions) Serialize() []string {
	args := []string{}
	if options.ExportToZip {
		args = append(args, "-f z")
	}
	if options.ConvertJeol {
		args = append(args, "-j")
	}
	if options.ConvertXes {
		args = append(args, "-x")
	}
	if options.ConvertQlw {
		args = append(args, "-q")
	}
	if options.ConvertMap {
		args = append(args, "-m")
	}
	if options.ConvertLine {
		args = append(args, "-l")
	}
	if options.ConvertQMap {
		args = append(args, "-k")
	}
	if options.Loose {
		args = append(args, "-y")
	}
	if options.Recover {
		args = append(args, "-r")
	}
	if options.Debug {
		args = append(args, "-d")
	}

	return args
}
