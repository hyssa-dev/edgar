package edgar

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func ParseCikAndDocID(url string) (string, string) {
	var s1 string
	var d1, d2, d3, d4 int
	_, _ = fmt.Sscanf(url, "/cgi-bin/viewer?action=view&cik=%d&accession_number=%d-%d-%d%s", &d1, &d2, &d3, &d4, &s1)
	cik := fmt.Sprintf("%d", d1)
	an := fmt.Sprintf("%010d%d%d", d2, d3, d4)
	return cik, an
}

/*
This is the parsing of query page where we get the list of filings of a given types
ex: https://www.sec.gov/cgi-bin/browse-edgar?CIK=AAPL&owner=exclude&action=getcompany&type=10-Q&count=1&dateb=
Assumptions of the parser:
- There is interactive data available and there is a button that allows the user to click it
- Since it is a link the tag will be a hyperlink with a button with the id=interactiveDataBtn
- The actual link is the href attribute in the "a" token just before the id attribute
*/
func QueryPageParser(page io.Reader, docType FilingType) map[string]string {

	filingInfo := make(map[string]string)

	z := html.NewTokenizer(page)

	_, data, err := ParseTableRow("query", z, true)
	for err == nil {
		//This check for filing type will drop AMEND filings
		if len(data) == 5 && data[0] == string(docType) {
			//Drop filings before 2010
			year := getYear(data[3])
			if year >= thresholdYear {
				filingInfo[data[3]] = data[1]
			}
		}
		_, data, err = ParseTableRow("query", z, true)
	}
	return filingInfo
}

func CikPageParser(page io.Reader) (string, error) {
	z := html.NewTokenizer(page)
	token := z.Token()
	for !(token.Data == "cik" && token.Type == html.StartTagToken) {
		tt := z.Next()
		if tt == html.ErrorToken {
			return "", errors.New("Could not find the CIK")
		}
		token = z.Token()
	}
	for !(token.Data == "cik" && token.Type == html.EndTagToken) {
		if token.Type == html.TextToken {
			str := strings.TrimSpace(token.String())
			if len(str) > 0 {
				return str, nil
			}
		}
		z.Next()
		token = z.Token()
	}
	return "", errors.New("Could not find the CIK")
}

/*
The filing page parser
- The top of the page has a list of reports.
- Get all the reports (link to all the reports) and put it in an array
- The Accordian on the side of the page identifies what each report is
- Get the text of the accordian and map the type of the report to the report
- Create a map of the report to report link
*/
func FilingPageParser(page io.Reader, _ FilingType) map[string][]Document {
	var filingLinks []string
	r := bufio.NewReader(page)
	s, e := r.ReadString('\n')

	for e == nil {
		//Get the number of reports available
		if strings.Contains(s, "var reports") {
			s1 := strings.Split(s, "(")
			s2 := strings.Split(s1[1], ")")
			cnt, _ := strconv.Atoi(s2[0])

			//cnt-1 because we skip the 'all' in the list
			for i := 0; i < cnt-1; i++ {
				if s, e = r.ReadString('\n'); e != nil {
					panic(e.Error())
				}
				s1 := strings.Split(s, " = ")
				s2 := strings.Split(s1[1], ";")
				s3 := strings.Trim(s2[0], "\"")
				s4 := strings.Split(s3, ".")
				s5 := s3
				//Sometimes the report is listed as an xml file??
				if s4[1] == "xml" {
					s5 = s4[0] + ".htm"
				}
				if !strings.Contains(s5, "htm") {
					panic("Dont know this type of report")
				}
				filingLinks = append(filingLinks, s5)
			}

			break
		}
		s, e = r.ReadString('\n')

	}

	docs := mapReports(page, filingLinks)
	return docs

}

func ParseTableData(z *html.Tokenizer, parseHref bool) (string, string) {
	token := z.Token()
	if token.Type != html.StartTagToken && token.Data != "td" {
		log.Fatal("Tokenizer passed incorrectly to parseTableData")
		return "", ""
	}

	for !(token.Data == "td" && token.Type == html.EndTagToken) {
		if token.Type == html.ErrorToken {
			break
		}

		if parseHref && token.Data == "a" && token.Type == html.StartTagToken {
			str, _ := ParseHyperLinkTag(z, token)
			// fmt.Println(way, str)
			if len(str) > 0 {
				return ParseKeyName(string(z.Buffered())), str
			}
		} else {
			if token.Type == html.TextToken {
				str := strings.TrimSpace(token.String())
				if len(str) > 0 {
					// fmt.Println("strings.TrimSpace", len(str), str)
					return "-1", str
				}
			}
		}
		//Going for the end of the td tag
		z.Next()
		token = z.Token()
	}
	return "", ""
}

func ParseTableRow(useCase string, z *html.Tokenizer, parseHref bool) ([]string, []string, error) {
	var (
		valueData []string
		keyData   []string
	)
	//Get the current token
	token := z.Token()
	//Check if this is really a table row
	for !(token.Type == html.StartTagToken && token.Data == "tr") {
		tt := z.Next()
		if tt == html.ErrorToken {
			return nil, nil, errors.New("Done with parsing")
		}
		token = z.Token()
	}
	//Till the end of the row collect data from each data block
	for !(token.Data == "tr" && token.Type == html.EndTagToken) {

		if token.Type == html.ErrorToken {
			return nil, nil, errors.New("Done with parsing")
		}
		if token.Data == "td" && token.Type == html.StartTagToken {
			parseFlag := parseHref
			//If the data is a number class just get the text = number
			for _, a := range token.Attr {
				if a.Key == "class" && (a.Val == "nump" || a.Val == "num") {
					parseFlag = false
				}
			}
			key, str := ParseTableData(z, parseFlag)

			// if key != "-1" && str != "" {
			// 	fmt.Println("-----------------------------------------------------------")
			// }

			// if str != "" {
			// 	fmt.Println(key, str)
			// }
			// if key != "-1" && str != "" {
			// 	fmt.Println()
			// 	fmt.Println("-----------------------------------------------------------")
			// }
			if len(str) > 0 {
				valueData = append(valueData, str)
				if key != "-1" {
					keyData = append(keyData, key)
				}
			}
		}
		z.Next()
		token = z.Token()
	}
	// fmt.Println("rowData", len(valueData), valueData)
	return keyData, valueData, nil
}

var reqHyperLinks = map[string]bool{
	"interactiveDataBtn": true,
}

func ParseHyperLinkTag(z *html.Tokenizer, token html.Token) (string, string) {
	// fmt.Println(token.String(), "---------------------------------------------------")
	// fmt.Println(string(z.Buffered()), z.Token(), string(z.Text()))
	// fmt.Println(token.String(), "---------------------------------------------------")
	// fmt.Scanln()
	var href string
	var onclick string
	var id string
	for _, a := range token.Attr {
		switch a.Key {
		case "id":
			id = a.Val
		case "href":
			href = a.Val
		case "onclick":
			onclick = a.Val
			if str, err := getFinDataXBRLTag(onclick); err == nil {
				return str, "onclick"
			}
		}
	}

	text := ""
	//Finish up the hyperlink
	for !(token.Data == "a" && token.Type == html.EndTagToken) {
		/*
			if token.Type == html.TextToken {
				str := strings.TrimSpace(token.String())
				if len(str) > 0 {
					text = str
				}
			}
		*/
		z.Next()
		token = z.Token()
	}

	if _, ok := reqHyperLinks[id]; ok {
		return href, "href"
	}

	return text, "text"
}

func ParseTableTitle(z *html.Tokenizer) []string {
	// fmt.Println("zzzzzzzzzz", string(z.Buffered()))
	var strs []string
	token := z.Token()

	if token.Type != html.StartTagToken && token.Data != "th" {
		log.Fatal("Tokenizer passed incorrectly to parseTableData")
		return strs
	}

	for !(token.Data == "th" && token.Type == html.EndTagToken) {
		if token.Type == html.ErrorToken {
			break
		}
		if token.Type == html.TextToken {
			str := strings.TrimSpace(token.String())
			if len(str) > 0 {
				strs = append(strs, str)
			}
		}
		//Going for the end of the td tag
		z.Next()
		token = z.Token()
	}
	return strs
}

func ParseTableHeading(z *html.Tokenizer) ([]string, error) {
	var retData []string
	//Get the current token
	token := z.Token()

	//Check if this is really a table row
	for !(token.Type == html.StartTagToken && token.Data == "tr") {
		tt := z.Next()
		if tt == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		token = z.Token()
	}

	//Till the end of the row collect data from each data block
	for !(token.Data == "tr" && token.Type == html.EndTagToken) {

		if token.Type == html.ErrorToken {
			return nil, errors.New("Done with parsing")
		}
		if token.Data == "th" && token.Type == html.StartTagToken {
			str := ParseTableTitle(z)
			if len(str) > 0 {
				retData = append(retData, str...)
			}
		}
		z.Next()
		token = z.Token()
	}

	return retData, nil
}

func ParseFilingScale(z *html.Tokenizer, docType string) (map[unitEntity]unit, []string, []string) {
	scales := make(map[unitEntity]unit)
	data, err := ParseTableHeading(z)
	z.Next()
	if err == nil {
		if len(data) > 0 {
			scales = filingScale(data, docType)
		}
	}
	periods, err := ParseTableHeading(z)
	if err == nil {
		if len(periods) > 0 {
			endDayes := 0
			for _, d := range data {
				if strings.Contains(d, "Months Ended") {
					endDayes += 1
				}
			}
			periods = periods[:len(data[endDayes:])]
		}
	}
	return scales, data, periods
}

/*
	This function takes any report filed under a company and looks
	for XBRL tags that Filing is interested in to gather and store.
	XBRL tag is mapped to a finDataType which is then used to lookup
	the passed in interface fields to see if there is a match and set
	that field
*/

func FinReportParser(page io.Reader, fr *FinancialReport, docType string) (*FinancialReport, error) {

	docValue := DocValues{
		Periods:  []string{},
		Section:  []Section{},
		EndDates: []string{},
		Scales:   map[unitEntity]unit{},
	}

	z := html.NewTokenizer(page)
	section := Section{
		Rows: make([]Row, 0),
	}

	scales, heading, periods := ParseFilingScale(z, docType)
	docValue.Scales = scales
	if len(heading) > 0 {
		fr.Title = strings.TrimSpace(heading[0])
		docValue.EndDates = heading[2:]
	}

	if len(periods) > 0 {
		docValue.Periods = periods
	}

	keys, values, err := ParseTableRow("fin", z, true)
	for err == nil {
		if len(keys) > 0 && strings.HasPrefix(keys[0], ">") {
			keys[0] = keys[0][1:]
			if len(docValue.Section) == 0 {
				section.Name = keys[0]
				section.Index = len(docValue.Section)
				section.Key = values[0]
			} else {
				docValue.Section = append(docValue.Section, section)
				section = Section{
					Name:  keys[0],
					Key:   values[0],
					Index: len(docValue.Section),
				}
			}
		}

		if len(values) > 0 {
			finType, continu := getFinDataTypeFromXBRLTag(values[0])
			if !continu {
				break
			}

			if len(values) > 1 {
				section.Rows = append(section.Rows, Row{
					Index: len(section.Rows),
					Name: func() string {
						if len(keys) > 0 {
							return keys[0]
						}

						return ""
					}(),
					Key:    values[0],
					Type:   finType,
					Values: values[1:],
				})
				for _, str := range values[1:] {
					if len(str) > 0 {
						if setData(fr, finType, str, scales, docType) == nil {
							break
						}
					}
				}
			}
		}
		keys, values, err = ParseTableRow("fin", z, true)
	}

	if len(section.Rows) > 0 {
		docValue.Section = append(docValue.Section, section)
	}
	fr.DocValues[docType] = append(fr.DocValues[docType], docValue)
	return fr, nil
}

// parseAllReports gets all the reports filed under a given account normalizeNumber
//
//nolint:unused // It is used in the parser
func ParseAllReports(cik string, an string) []int {
	var reports []int
	url := "https://www.sec.gov/Archives/edgar/data/" + cik + "/" + an + "/"
	page := getPage(url)
	z := html.NewTokenizer(page)
	_, data, err := ParseTableRow("all", z, false)
	for err == nil {
		var num int
		if len(data) > 0 && strings.Contains(data[0], "R") {
			_, err = fmt.Sscanf(data[0], "R%d.htm", &num)
			if err == nil {
				reports = append(reports, num)
			}
		}
		_, data, err = ParseTableRow("all", z, false)
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i] < reports[j]
	})
	return reports
}

func ParseMappedReports(typeDocs map[string][]Document, filingType FilingType) (*FinancialReport, error) {
	var wg sync.WaitGroup

	fr := newFinancialReport(filingType)
	for docType, docs := range typeDocs {
		fr = fr.newDocValues(docType)
		for _, doc := range docs {
			url := baseURL + doc.URL
			wg.Add(1)
			go func(url string, fr *FinancialReport, t string) {
				defer wg.Done()
				page := getPage(url)
				if page != nil {
					_, _ = FinReportParser(page, fr, docType)
				}
			}(url, fr, docType)
		}
	}
	wg.Wait()
	return fr, nil //validateFinancialReport(fr)
}

func ParseKeyName(data string) string {
	if strings.HasPrefix(data, "<strong>") {
		idx := strings.Index(data, "</strong>")
		if idx > 0 {
			return data[7:idx]
		}
	} else {
		idx := strings.Index(data, "<")
		if idx > 0 {
			return data[:idx]
		}
	}

	return ""
}
