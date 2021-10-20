package imageset

type SyncSet struct {
	Names    []string `json:"names"`
	Prefixes []string `json:"prefixes"`
}

type ImageSet struct {
	Sync    *SyncSet `json:"sync"`
	Install []string `json:"install"`
}
