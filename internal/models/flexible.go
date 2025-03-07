package models

type FlexExecuteOptions struct {
	MaxCount        int
	NoAnalyze       bool
	ColumnBlacklist []string
	NoBlacklist     bool
}

type FlexSearchResponse struct {
	Headers    []string       `json:"headers"`
	ResultList [][]string     `json:"resultList"`
	Exception  *FlexException `json:"exception"`
}

type FlexException struct {
	Message string `json:"message"`
}
