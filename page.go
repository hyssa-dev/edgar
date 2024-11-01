package edgar

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	baseURL              = "https://www.sec.gov/"
	cikURL               = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&output=xml&CIK=%s"
	queryURL             = "cgi-bin/browse-edgar?action=getcompany&CIK=%s&type=%s&dateb=&owner=exclude&count=10"
	searchURL            = baseURL + queryURL
	acceptEncodingHeader = "Accept-Encoding"
	acceptEncodingValue  = "deflate, br, zstd"
	userAgentHeader      = "User-Agent"
	userAgentValue       = "CompanyName company@vgmail.com"
	hostHeader           = "Host"
	hostValue            = "www.sec.gov"
	acceptHeader         = "Accept"
	acceptValue          = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
)

func createQueryURL(symbol string, docType FilingType) string {
	return fmt.Sprintf(searchURL, symbol, docType)
}

func getPage(url string) io.ReadCloser {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Query to SEC page ", url, "failed: ", err)
		return nil
	}

	req.Header.Add(userAgentHeader, userAgentValue)
	req.Header.Add(acceptEncodingHeader, acceptEncodingValue)
	req.Header.Add(hostHeader, hostValue)
	req.Header.Add(acceptHeader, acceptValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Query to SEC page ", url, "failed: ", err)
		return nil
	}
	return resp.Body
}

func getCompanyCIK(ticker string) string {
	url := fmt.Sprintf(cikURL, ticker)
	r := getPage(url)
	if r != nil {
		if cik, err := CikPageParser(r); err == nil {
			return cik
		}
	}
	return ""
}

// getFilingLinks gets the links for filings of a given type of filing 10K/10Q..
func getFilingLinks(ticker string, fileType FilingType) map[string]string {
	url := createQueryURL(ticker, fileType)
	resp := getPage(url)
	if resp == nil {
		log.Println("No response on the query for docs")
		return nil
	}
	defer resp.Close()
	return QueryPageParser(resp, fileType)

}

// Get all the docs pages based on the filing type
// Returns a map:
// key=Document type ex.Cash flow statement
// Value = link to that that sheet
func getFilingDocs(url string, fileType FilingType) map[string][]Document {
	url = baseURL + url
	resp := getPage(url)
	if resp == nil {
		return nil
	}
	defer resp.Close()
	return FilingPageParser(resp, fileType)
}

// getFinancialData gets the data from all the filing docs and places it in
// a financial report
func getFinancialData(url string, fileType FilingType) (*Report, error) {
	docs := getFilingDocs(url, fileType)
	report := Report{
		Documents: docs,
	}

	financial, err := ParseMappedReports(docs, fileType)
	if err != nil {
		return nil, err
	}
	report.FinData = financial

	return &report, nil
}
