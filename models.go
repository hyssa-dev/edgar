package edgar

import "sync"

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
	URL        string
	IsRequired bool
}

type Company struct {
	sync.Mutex
	Company     string `json:"Company"`
	cik         string
	FilingLinks map[FilingType]map[string]string  `json:"-"`
	Reports     map[FilingType]map[string]*Report `json:"Financial Reports"`
}

type Report struct {
	Company string           `json:"Company"`
	Date    Timestamp        `json:"Report date"`
	FinData *FinancialReport `json:"Financial Data"`

	// DocumentType/Document
	Documents map[string][]Document `json:"Documents"`
}

type DocValues struct {
	Periods  []string
	Section  []Section
	EndDates []string
	Scales   map[unitEntity]unit
}

type Section struct {
	Name  string
	Key   string
	Index int
	Unit  int
	Rows  []Row
}

type Row struct {
	Index  int
	Name   string
	Key    string
	Type   string
	Values []string
}
