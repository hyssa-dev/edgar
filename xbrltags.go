package edgar

import (
	"strings"
)

var (
	// A Map of XBRL tags to financial data type
	// This map contains the corresponding GAAP tag and a version of the tag
	// without the GAAP keyword in case the company has only file non-gaap
	xbrlTags = map[string]string{
		//Balance Sheet info
		"defref_us-gaap_StockholdersEquity":                            "Total Shareholder Equity",
		"StockholdersEquity":                                           "Total Shareholder Equity",
		"defref_us-gaap_RetainedEarningsAccumulatedDeficit":            "Retained Earnings",
		"RetainedEarningsAccumulatedDeficit":                           "Retained Earnings",
		"defref_us-gaap_LiabilitiesCurrent":                            "Current Liabilities",
		"LiabilitiesCurrent":                                           "Current Liabilities",
		"defref_us-gaap_AssetsCurrent":                                 "Current Assets",
		"AssetsCurrent":                                                "Current Assets",
		"defref_us-gaap_Assets":                                        "Total Assets",
		"Assets":                                                       "Total Assets",
		"defref_us-gaap_Liabilities":                                   "Total Liabilities",
		"Liabilities":                                                  "Total Liabilities",
		"defref_us-gaap_CashAndCashEquivalentsAtCarryingValue":         "Cash",
		"CashAndCashEquivalentsAtCarryingValue":                        "Cash",
		"defref_us-gaap_Goodwill":                                      "Goodwill",
		"Goodwill":                                                     "Goodwill",
		"defref_us-gaap_IntangibleAssetsNetExcludingGoodwill":          "Intangibles",
		"IntangibleAssetsNetExcludingGoodwill":                         "Intangibles",
		"defref_us-gaap_LongTermDebtNoncurrent":                        "Long-Term debt",
		"LongTermDebtNoncurrent":                                       "Long-Term debt",
		"defref_us-gaap_LongTermDebtAndCapitalLeaseObligations":        "Long-Term debt",
		"LongTermDebtAndCapitalLeaseObligations":                       "Long-Term debt",
		"defref_us-gaap_ShortTermBorrowings":                           "Short-Term debt",
		"ShortTermBorrowings":                                          "Short-Term debt",
		"defref_us-gaap_DebtCurrent":                                   "Short-Term debt",
		"DebtCurrent":                                                  "Short-Term debt",
		"defref_us-gaap_LongTermDebtAndCapitalLeaseObligationsCurrent": "Short-Term debt",
		"LongTermDebtAndCapitalLeaseObligationsCurrent":                "Short-Term debt",
		"defref_us-gaap_DeferredRevenueCurrent":                        "Deferred revenue",
		"DeferredRevenueCurrent":                                       "Deferred revenue",
		"defref_us-gaap_RetainedEarningsAccumulatedDeficitAndAccumulatedOtherComprehensiveIncomeLossNetOfTax": "Retained Earnings",
		"RetainedEarningsAccumulatedDeficitAndAccumulatedOtherComprehensiveIncomeLossNetOfTax":                "Retained Earnings",

		//Operations Sheet info
		"defref_us-gaap_SalesRevenueNet": "Revenue",
		"SalesRevenueNet":                "Revenue",
		"defref_us-gaap_Revenues":        "Revenue",
		"Revenues":                       "Revenue",
		"defref_us-gaap_RevenueFromContractWithCustomerExcludingAssessedTax":            "Revenue",
		"RevenueFromContractWithCustomerExcludingAssessedTax":                           "Revenue",
		"defref_us-gaap_CostOfRevenue":                                                  "Cost Of Revenue",
		"defref_us-gaap_CostOfGoodsAndServicesSold":                                     "Cost Of Revenue",
		"CostOfGoodsAndServicesSold":                                                    "Cost Of Revenue",
		"defref_us-gaap_CostOfPurchasedOilAndGas":                                       "Cost Of Revenue",
		"CostOfPurchasedOilAndGas":                                                      "Cost Of Revenue",
		"defref_us-gaap_CostOfGoodsSold":                                                "Cost Of Revenue",
		"CostOfGoodsSold":                                                               "Cost Of Revenue",
		"defref_us-gaap_CostOfGoodsSoldExcludingAmortizationOfAcquiredIntangibleAssets": "Cost Of Revenue",
		"CostOfGoodsSoldExcludingAmortizationOfAcquiredIntangibleAssets":                "Cost Of Revenue",
		"defref_us-gaap_GrossProfit":                                                    "Gross Margin",
		"GrossProfit":                                                                   "Gross Margin",
		"defref_us-gaap_OperatingExpenses":                                              "Operational Expense",
		"OperatingExpenses":                                                             "Operational Expense",
		"defref_us-gaap_CostsAndExpenses":                                               "Operational Expense",
		"CostsAndExpenses":                                                              "Operational Expense",
		"defref_us-gaap_OtherCostAndExpenseOperating":                                   "Operational Expense",
		"OtherCostAndExpenseOperating":                                                  "Operational Expense",
		"defref_us-gaap_OperatingIncomeLoss":                                            "Operational Income",
		"OperatingIncomeLoss":                                                           "Operational Income",
		"defref_us-gaap_IncomeLossIncludingPortionAttributableToNoncontrollingInterest": "Operational Income",
		"defref_us-gaap_IncomeLossFromContinuingOperationsIncludingPortionAttributableToNoncontrollingInterest":                      "Operational Income",
		"IncomeLossFromContinuingOperationsIncludingPortionAttributableToNoncontrollingInterest":                                     "Operational Income",
		"IncomeLossIncludingPortionAttributableToNoncontrollingInterest":                                                             "Operational Income",
		"defref_us-gaap_IncomeLossFromContinuingOperationsBeforeIncomeTaxesMinorityInterestAndIncomeLossFromEquityMethodInvestments": "Operational Income",
		"IncomeLossFromContinuingOperationsBeforeIncomeTaxesMinorityInterestAndIncomeLossFromEquityMethodInvestments":                "Operational Income",
		"defref_us-gaap_NetIncomeLoss": "Net Income",
		"NetIncomeLoss":                "Net Income",
		"defref_us-gaap_ProfitLoss":    "Net Income",
		"ProfitLoss":                   "Net Income",
		"defref_us-gaap_NetIncomeLossAvailableToCommonStockholdersBasic": "Net Income",
		"NetIncomeLossAvailableToCommonStockholdersBasic":                "Net Income",
		"defref_us-gaap_WeightedAverageNumberOfDilutedSharesOutstanding": "Weighted Average Share Count",
		"WeightedAverageNumberOfDilutedSharesOutstanding":                "Weighted Average Share Count",
		"defref_us-gaap_CommonStockDividendsPerShareDeclared":            "Dividend Per Share",
		"CommonStockDividendsPerShareDeclared":                           "Dividend Per Share",

		"defref_us-gaap_IncomeLossFromContinuingOperationsBeforeIncomeTaxesExtraordinaryItemsNoncontrollingInterest": "Operational Income",
		"IncomeLossFromContinuingOperationsBeforeIncomeTaxesExtraordinaryItemsNoncontrollingInterest":                "Operational Income",

		//Cash Flow Sheet info
		"defref_us-gaap_NetCashProvidedByUsedInOperatingActivities":                     "Operating Cash Flow",
		"NetCashProvidedByUsedInOperatingActivities":                                    "Operating Cash Flow",
		"defref_us-gaap_NetCashProvidedByUsedInOperatingActivitiesContinuingOperations": "Operating Cash Flow",
		"NetCashProvidedByUsedInOperatingActivitiesContinuingOperations":                "Operating Cash Flow",
		"defref_us-gaap_PaymentsToAcquirePropertyPlantAndEquipment":                     "Capital Expenditure",
		"PaymentsToAcquirePropertyPlantAndEquipment":                                    "Capital Expenditure",
		"defref_us-gaap_PaymentsToAcquireProductiveAssets":                              "Capital Expenditure",
		"PaymentsToAcquireProductiveAssets":                                             "Capital Expenditure",
		"defref_us-gaap_CapitalExpendituresAndInvestments":                              "Capital Expenditure",
		"CapitalExpendituresAndInvestments":                                             "Capital Expenditure",
		"defref_us-gaap_PaymentsOfDividends":                                            "Dividends paid",
		"PaymentsOfDividends":                                                           "Dividends paid",
		"defref_us-gaap_PaymentsOfDividendsCommonStock":                                 "Dividends paid",
		"PaymentsOfDividendsCommonStock":                                                "Dividends paid",
		"defref_us-gaap_InterestPaidNet":                                                "Interest paid",
		"InterestPaidNet":                                                               "Interest paid",
		"defref_us-gaap_InterestAndDebtExpense":                                         "Interest paid",
		"InterestAndDebtExpense":                                                        "Interest paid",
		"defref_us-gaap_InterestIncomeExpenseNet":                                       "Interest paid",
		"InterestIncomeExpenseNet":                                                      "Interest paid",
		//Entity sheet information
		"defref_dei_EntityCommonStockSharesOutstanding": "Shares Outstanding",
		"EntityCommonStockSharesOutstanding":            "Shares Outstanding",
	}
)

func getFinDataTypeFromXBRLTag(key string) string {
	data, ok := xbrlTags[key]
	if !ok {

		// Now look for non-gaap filing
		// defref_us-gaap_XXX could be filed company specific
		// as defref_msft_XXX
		splits := strings.Split(key, "_")
		if len(splits) == 3 {
			data, ok = xbrlTags[splits[2]]
			if ok {
				return data
			}
		}
		return unknownDataType
	}
	return data
}
