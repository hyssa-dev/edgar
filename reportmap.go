package edgar

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var (
	unknownMenuCat  = "Unknown"
	unknowdDocType  = "Unknown"
	unknownDataType = "Unknown"
)

// This function returns the menu category of the document
func getMenuCategory(menuCategories []MenuCategory, data string) string {
	if len(menuCategories) == 0 {
		menuCategories = MenuCategories
	}

	data = strings.ToLower(data)
	// fmt.Println("data", data, "menuCategories", menuCategories)
	for _, category := range menuCategories {
		// fmt.Println("category", category, "data", data, "contains", containsAllElements(data, category.Keys), "nonContains", nonContainsAllElements(data, category.NotKeys))
		if category.Condition == "all" {
			if containsAllElements(data, category.Keys) && nonContainsAllElements(data, category.NotKeys) {
				return category.Name
			}
		} else {
			if containsAnyElement(data, category.Keys) && nonContainsAllElements(data, category.NotKeys) {
				return category.Name
			}
		}
	}

	return unknownMenuCat
}

// This function returns the filing type of the document(by the menu category)
func lookupDocType(data string, menuCategory string, categoryDocs map[string][]Document) (Document, error) {
	data = strings.ToUpper(data)

	if len(categoryDocs) == 0 {
		categoryDocs = CategoryDocs
	}
	docs := categoryDocs[menuCategory]
	for _, doc := range docs {
		if doc.Condition == "all" {
			if containsAllElements(data, doc.Keys) && nonContainsAllElements(data, doc.NotKeys) {
				return doc, nil
			}
		} else {
			if containsAnyElement(data, doc.Keys) && nonContainsAllElements(data, doc.NotKeys) {
				return doc, nil
			}
		}
	}

	return Document{}, errors.New("not found")
}

// NEED TO REILIZE THIS FUNCTION: TOOD
func getMissingDocs(urlByDocType map[string][]Document, requiredDocs []Document) string {
	if len(requiredDocs) == 0 {
		requiredDocs = RequiredDocs
	}
	// fmt.Println("urlByDocType", urlByDocType)
	// fmt.Println()
	// fmt.Println("requiredDocs", requiredDocs)
	if len(urlByDocType) >= len(requiredDocs) {
		return ""
	}
	var diff []string
	for _, doc := range requiredDocs {
		if _, ok := urlByDocType[doc.Type]; !ok {
			switch doc.Type {
			case "Operations":
				if _, ok := urlByDocType["Income"]; ok {
					continue
				}
			case "Income":
				if _, ok := urlByDocType["Operations"]; ok {
					continue
				}
			}
			// fmt.Println(doc.Type)
			diff = append(diff, doc.Type)
		}
	}
	if len(diff) == 0 {
		return ""
	}

	var ret string
	ret = "[ "
	for _, val := range diff {
		ret = ret + " " + string(val)
	}
	ret += " ]"
	return ret
}

func mapReports(page io.Reader, filingLinks []string) map[string][]Document {

	menuCategory := unknownMenuCat

	urlByDocType := make(map[string][]Document)

	z := html.NewTokenizer(page)
	tt := z.Next()
loop:
	for tt != html.ErrorToken {
		token := z.Token()
		if token.Data == "a" {
			for _, a := range token.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "loadReport") {
					strs := strings.Split(a.Val, "loadReport")
					strs[1] = strings.Trim(strs[1], ";")
					reportNum, _ := strconv.Atoi(strings.Trim(strs[1], "()"))
					tt = z.Next() //Contains the text that describes the report
					if tt != html.TextToken {
						break
					}
					token = z.Token()
					doc, errDoc := lookupDocType(token.String(), menuCategory, map[string][]Document{})
					if errDoc != nil {
						doc.Type = unknowdDocType
					}
					doc.URL = filingLinks[reportNum-1]
					urlByDocType[doc.Type] = append(urlByDocType[doc.Type], doc)
				} else if a.Key == "id" && strings.Contains(a.Val, "menu_cat") {
					// Set the menu level
					for !(token.Data == "a" && token.Type == html.EndTagToken) {
						if token.Type == html.TextToken {
							str := strings.TrimSpace(token.String())
							menuCategory = getMenuCategory([]MenuCategory{}, str)
						}
						z.Next()
						token = z.Token()
					}
					if menuCategory == unknownMenuCat {
						//Gone too far. Menu category 4 is beyond notes of financial statements.
						//Stop parsing
						break loop
					}
				}
			}
		}
		tt = z.Next()
	}
	ret := getMissingDocs(urlByDocType, []Document{})
	if ret != "" {
		log.Println("Did not find the following filing documents: " + ret)
	}
	return urlByDocType
}
