package edgar

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"
)

func (f Report) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Filing data")
	}
	return string(data)
}

func (f *Report) Ticker() string {
	return f.Company
}

func (f *Report) FiledOn() time.Time {
	return time.Time(f.Date)
}

func (f *Report) filingErrorString() string {
	return "Filing information has not been collected for " + getDateString(f.FiledOn()) + " "
}

func (f *Report) Type() (FilingType, error) {
	if f.FinData != nil {
		return f.FinData.DocType, nil
	}
	return "", errors.New(f.filingErrorString())
}

func (f *Report) ShareCount() (float64, error) {
	if f.FinData != nil && f.FinData.Entity != nil {
		if isCollectedDataSet(f.FinData.Entity, "ShareCount") {
			return f.FinData.Entity.ShareCount, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Share Count")
}

func (f *Report) Revenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "Revenue") {
			return f.FinData.Ops.Revenue, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Revenue")
}

func (f *Report) CostOfRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "CostOfSales") {
			return f.FinData.Ops.CostOfSales, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Cost of Revenue")
}

func (f *Report) GrossMargin() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "GrossMargin") {
			return f.FinData.Ops.GrossMargin, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Gross Margin")
}

func (f *Report) OperatingIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "OpIncome") {
			return f.FinData.Ops.OpIncome, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Income")
}

func (f *Report) OperatingExpense() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "OpExpense") {
			return f.FinData.Ops.OpExpense, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Expense")
}

func (f *Report) NetIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "NetIncome") {
			return f.FinData.Ops.NetIncome, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Net Income")
}

func (f *Report) TotalEquity() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Equity") {
			return f.FinData.Bs.Equity, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Total Equity")
}

func (f *Report) ShortTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "SDebt") {
			return f.FinData.Bs.SDebt, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Short Term Debt")
}

func (f *Report) LongTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "LDebt") {
			return f.FinData.Bs.LDebt, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Long Term Debt")
}

func (f *Report) CurrentLiabilities() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "CLiab") {
			return f.FinData.Bs.CLiab, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Current Liabilities")
}

func (f *Report) CurrentAssets() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "CAssets") {
			return f.FinData.Bs.CAssets, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Current Assets")
}

func (f *Report) DeferredRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Deferred") {
			return f.FinData.Bs.Deferred, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Deferred Revenue")
}

func (f *Report) RetainedEarnings() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Retained") {
			return f.FinData.Bs.Retained, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Retained Earnings")
}

func (f *Report) OperatingCashFlow() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "OpCashFlow") {
			return f.FinData.Cf.OpCashFlow, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Cash Flow")
}

func (f *Report) CapitalExpenditure() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "CapEx") {
			return f.FinData.Cf.CapEx, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Capital Expenditur")
}

func (f *Report) Dividend() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "Dividends") {
			// Dividend is recorded as an expense and is -ve. Hence reversing sign
			return f.FinData.Cf.Dividends * -1, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Dividend")
}

func (f *Report) WAShares() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "WAShares") {
			return f.FinData.Ops.WAShares, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Weighted Average Shares")
}

func (f *Report) DividendPerShare() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "Dps") {
			return f.FinData.Ops.Dps, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Dividend Per Share")
}

func (f *Report) Interest() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "Interest") {
			return f.FinData.Cf.Interest, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Interest paid")
}

func (f *Report) Cash() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Cash") {
			return f.FinData.Bs.Cash, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Cash")
}

func (f *Report) Securities() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Securities") {
			return f.FinData.Bs.Securities, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Securities")
}

func (f *Report) Goodwill() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Goodwill") {
			return f.FinData.Bs.Goodwill, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Goodwill")
}

func (f *Report) Intangibles() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Intangibles") {
			return f.FinData.Bs.Intangibles, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Intangibles")
}

func (f *Report) Assets() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Assets") {
			return f.FinData.Bs.Assets, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Total Assets")
}

func (f *Report) Liabilities() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Liab") {
			return f.FinData.Bs.Liab, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Total Liabilities")
}

func (f *Report) CollectedData() []string {

	eval := func(data interface{}) []string {
		var ret []string
		if data != nil {
			t := reflect.TypeOf(data)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			for i := 0; i < t.NumField(); i++ {
				if isCollectedDataSet(data, t.Field(i).Name) {
					ret = append(ret, t.Field(i).Name)
				}
			}
		}
		return ret
	}
	ret := eval(f.FinData.Entity)
	ret = append(ret, eval(f.FinData.Bs)...)
	ret = append(ret, eval(f.FinData.Cf)...)
	ret = append(ret, eval(f.FinData.Ops)...)

	return ret
}
