package deps

type ImageSet struct {
	Names    []string `json:"names"`
	Prefixes []string `json:"prefixes"`
	Exclude  []string `json:"exclude"`
	Patches  []string `json:"patches"`
}

type Version struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
	Uri  string `json:"uri"`
}

type Versions struct {
	Arch       string     `json:"arch"`
	Components []*Version `json:"components"`
	Images     []*Version `json:"images"`
}

type TmplData struct {
	Tag  string
	Arch string
}
