package edgar

type MenuCategory struct {
	Name    string
	Keys    []string
	NotKeys []string
}

type Document struct {
	Category   string
	Name       string
	Type       string
	Keys       []string
	NotKeys    []string
	IsRequired bool
}

type Documents struct {
	Docs []Document
}
