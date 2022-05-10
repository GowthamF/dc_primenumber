package models

type PrimeNumbers struct {
	NumberToCheck int32 `json:"numberToCheck"`
	StartRange    int32 `json:"startRange"`
	EndRange      int32 `json:"endRange"`
}

type PrimeNumbersValidationMessage struct {
	NumberToCheck   int32 `json:"numberToCheck"`
	IsPrimeNumber   *bool `json:"isPrimeNumber"`
	DivisibleNumber int32 `json:"divisibleNumber"`
}
