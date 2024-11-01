package edgar

var MenuCategories = []MenuCategory{
	{
		Name: "Cover",
		Keys: []string{"cover"},
	},
	{
		Name: "Financial statements",
		Keys: []string{"financial", "statement"},
	},
	{
		Name: "Notes to Financial statements",
		Keys: []string{"note", "financial", "statement"},
	},
	{
		Name: "Notes Tables",
		Keys: []string{"table"},
	},
	{
		Name: "Notes Details",
		Keys: []string{"detail"},
	},
	{
		Name: "Accounting Policies",
		Keys: []string{"polic"},
	},
}

var CategoryDocs = map[string][]Document{
	"Cover": {
		{
			Category:   "Cover",
			Name:       "Cover Page",
			Type:       "Entity Info",
			Keys:       []string{"DOCUMENT", "ENTITY"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Cover",
			Name:       "Cover Page",
			Type:       "Entity Info",
			Keys:       []string{"COVER"},
			NotKeys:    []string{},
			IsRequired: false,
		},
	},
	"Financial statements": {
		{
			Category:   "Financial statements",
			Name:       "Balance sheet",
			Type:       "Assets",
			Keys:       []string{"BALANCE SHEET"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Balance sheet",
			Type:       "Assets",
			Keys:       []string{"FINANCIAL POSITION"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Operations statement",
			Type:       "Operations",
			Keys:       []string{"OPERATIONS"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Income statement",
			Type:       "Income",
			Keys:       []string{"INCOME"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Income statement",
			Type:       "Income",
			Keys:       []string{"INCOME", "EARNINGS"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Cash flow statement",
			Type:       "Cash Flow",
			Keys:       []string{"CASH FLOWS"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Financial statements",
			Name:       "Parenthetical statement",
			Type:       "Parenthetical",
			Keys:       []string{"PARENTHETICAL"},
			NotKeys:    []string{},
			IsRequired: false,
		},
	},
	"Notes to Financial statements": {
		{
			Category:   "Notes to Financial statements",
			Name:       "Notes - Notes on EPS",
			Type:       "Notes on EPS",
			Keys:       []string{"EARNINGS", "SHARE"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Notes to Financial statements",
			Name:       "Notes - Notes on Equity",
			Type:       "Notes on Equity",
			Keys:       []string{"SHAREHOLDER", "EQUITY"},
			NotKeys:    []string{},
			IsRequired: false,
		},
		{
			Category:   "Notes to Financial statements",
			Name:       "Notes - Notes on Debt",
			Type:       "Notes on Debt",
			Keys:       []string{"DEBT"},
			NotKeys:    []string{},
			IsRequired: false,
		},
	},
	"Notes Tables": {{
		Category:   "Notes Tables",
		Name:       "Notes - Notes on EPS",
		Type:       "Tables",
		Keys:       []string{"(TABLES)"},
		NotKeys:    []string{},
		IsRequired: false,
	},
		{
			Category:   "Notes Tables",
			Name:       "Accounting Policies (Tables)",
			Type:       "PoliciesTable",
			Keys:       []string{"Accounting Policies (Tables)"},
			NotKeys:    []string{""},
			IsRequired: false,
		},
	},
	"Notes Details": {
		{
			Category:   "}Notes Details",
			Name:       "Notes Details",
			Type:       "Details",
			Keys:       []string{"(DETAILS)", "(DETAIL)"},
			NotKeys:    []string{""},
			IsRequired: false,
		},
	},

	"Accounting Policies": {
		{
			Category:   "Accounting Policies",
			Name:       "Accounting Policies",
			Type:       "Policies",
			Keys:       []string{"(POLICIES)"},
			NotKeys:    []string{""},
			IsRequired: false,
		},
	},
}

var RequiredDocs = []Document{
	{
		Category:   "",
		Name:       "Operations",
		Type:       "Operations",
		Keys:       []string{},
		NotKeys:    []string{},
		IsRequired: true,
	},
	{
		Name:       "Income",
		Type:       "Income",
		IsRequired: true,
	},
	{
		Name:       "Assets",
		Type:       "Assets",
		IsRequired: true,
	},
	{
		Name:       "Cash Flow",
		Type:       "Cash Flow",
		IsRequired: true,
	},
	{
		Name:       "Entity Info",
		Type:       "Entity Info",
		IsRequired: true,
	},
}
