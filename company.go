package edgar

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"sort"
	"sync"
	"time"
)

func (c *Company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	return string(data)
}

func newCompany(ticker string) *Company {
	return &Company{
		Company:     ticker,
		cik:         getCompanyCIK(ticker),
		FilingLinks: make(map[FilingType]map[string]string),
		Reports:     make(map[FilingType]map[string]*Report),
	}
}

func (c *Company) Ticker() string {
	return c.Company
}

func (c *Company) Filing(fileType FilingType, ts time.Time) (Filing, error) {
	file, ok := c.getReport(fileType, ts)
	if !ok {
		link, ok1 := c.getFilingLink(fileType, ts)
		if !ok1 {
			log.Println(c.AvailableFilings(fileType))
			return nil, errors.New("No filing available for given date " + getDateString(ts))
		}
		file = new(Report)

		resp, err := getFinancialData(link, fileType)
		file.Documents = resp.Documents
		file.FinData = resp.FinData
		if file.FinData != nil {
			file.Date = Timestamp(ts)
			file.Company = c.Ticker()
			c.AddReport(file)
			if err != nil {
				log.Println(file.Company + "-Filed on: " + getDateString(ts) + ":" + err.Error())
			}
			return file, nil
		}
		return nil, err
	}
	return file, nil
}

// Get multiple filings in parallel
func (c *Company) Filings(fileType FilingType, ts ...time.Time) ([]Filing, error) {
	var wg sync.WaitGroup
	var ret []Filing
	var retErrors []error
	var m sync.Mutex
	for _, t := range ts {
		wg.Add(1)
		go func(filed time.Time) {
			defer wg.Done()
			file, err := c.Filing(fileType, filed)
			m.Lock()
			if err == nil {
				ret = append(ret, file)
			} else {
				err = errors.New(getDateString(filed) + ":" + err.Error())
				retErrors = append(retErrors, err)
			}
			m.Unlock()
		}(t)
	}
	wg.Wait()
	if len(ts) != len(ret) && len(retErrors) > 0 {
		errString := "Failed to retrieve some filings: \n"
		for _, e := range retErrors {
			errString = errString + e.Error() + "\n"
		}
		return ret, errors.New(errString)
	}
	return ret, nil
}

func (c *Company) AddReport(file *Report) {
	t, err := file.Type()
	if err != nil {
		log.Fatal("Adding invalid report")
		return
	}
	c.Lock()
	defer c.Unlock()
	if c.Reports[t] == nil {
		c.Reports[t] = make(map[string]*Report)
	}
	c.Reports[t][file.Date.String()] = file
}

func (c *Company) getReport(fileType FilingType, ts time.Time) (*Report, bool) {
	c.Lock()
	defer c.Unlock()
	file, ok := c.Reports[fileType][getDateString(ts)]
	return file, ok
}

func (c *Company) AvailableFilings(filingType FilingType) []time.Time {
	var d []time.Time
	c.Lock()
	links := c.FilingLinks[filingType]
	for key := range links {
		d = append(d, time.Time(getDate(key)))
	}
	c.Unlock()
	sort.Slice(d, func(i, j int) bool {
		return d[i].After(d[j])
	})
	return d
}

func (c *Company) CIK() string {
	return c.cik
}

func (c *Company) getFilingLink(fileType FilingType, ts time.Time) (string, bool) {
	c.Lock()
	defer c.Unlock()
	link, ok := c.FilingLinks[fileType][getDateString(ts)]
	return link, ok
}

func (c *Company) addFilingLinks(fileType FilingType, files map[string]string) {
	c.Lock()
	defer c.Unlock()
	c.FilingLinks[fileType] = files
}

// Save the Company folder into the writer in JSON format
func (c *Company) SaveFolder(w io.Writer) error {
	_, err := w.Write([]byte(c.String()))
	if err != nil {
		log.Println("Failed to save data")
		return err
	}
	return nil
}
