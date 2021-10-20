package deps

type ImageSet struct {
	Names    []string `json:"names"`
	Prefixes []string `json:"prefixes"`
	Exclude  []string `json:"exclude"`
}

type Version struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type Versions struct {
	Arch       string     `json:"arch"`
	Components []*Version `json:"components"`
	Images     []*Version `json:"images"`
}
