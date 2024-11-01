package edgar

import (
	"strings"
)

func getFinDataTypeFromXBRLTag(key string) (string, bool) {
	data, ok := XBRLTags[key]
	if !ok {

		// Now look for non-gaap filing
		// defref_us-gaap_XXX could be filed company specific
		// as defref_msft_XXX
		splits := strings.Split(key, "_")
		if len(splits) == 3 {
			data, ok = XBRLTags[splits[2]]
			if ok {
				return data, false
			}
		}

		if RestrictedTags[key] {
			return "", true
		}

		return unknownDataType, false
	}
	return data, false
}
