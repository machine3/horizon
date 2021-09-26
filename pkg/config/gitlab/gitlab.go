package gitlab

type Config struct {
	Application          *Gitlab `yaml:"application"`
	ApplicationInstances *Gitlab `yaml:"applicationInstance"`
}

type Gitlab struct {
	GitlabName string  `yaml:"gitlabName"`
	Parent     *Parent `yaml:"parent"`
}

type Parent struct {
	Path string `yaml:"path"`
	ID   int    `yaml:"id"`
}
