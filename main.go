package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const AMFII_NAV_URL = "https://portal.amfiindia.com/DownloadNAVHistoryReport_Po.aspx"

type MutualFund struct {
	SchemeCode          string
	SchemeName          string
	ISINDivPayout       string
	ISINDivReinvestment string
	NetAssetValue       string
	RepurchasePrice     string
	SalePrice           string
	Date                string
}

type FundHouse struct {
	Name  string
	Funds []MutualFund
}

type SchemeCategory struct {
	Name       string
	FundHouses []FundHouse
}

func main() {
	fromDate := time.Now().AddDate(0, 0, -1).Format("02-Jan-2006")
	url := AMFII_NAV_URL + fmt.Sprintf("?frmdt=%s", fromDate)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("error fetching data from upstream API: %w", err)
	}
	defer resp.Body.Close()

	categories, err := readCSV(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	displayData(categories)
}

// readCSV opens the CSV file and processes its records into SchemeCategory structures.
// It returns a slice of SchemeCategory and an error if encountered.
func readCSV(data io.ReadCloser) ([]SchemeCategory, error) {
	reader := csv.NewReader(data)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	var categories []SchemeCategory
	var currentCategory *SchemeCategory
	var currentFundHouse *FundHouse

	// Read and process each record from the CSV
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip empty lines
		if skipEmptyLines(record) {
			continue
		}

		// Handle the record based on its length (1 for category/fund house, 8 for fund data)
		parseCategoryOrFundHouse(record, &categories, &currentCategory, &currentFundHouse)
		if len(record) == 8 {
			parseFund(record, currentFundHouse)
		}
	}

	// Finalize any remaining category and fund house after reading all records
	finalizeCategoryAndFundHouse(&categories, currentCategory, currentFundHouse)

	return categories, nil
}

// skipEmptyLines checks if a record contains only an empty string and returns true if so.
// This helps in skipping empty lines in the CSV file.
func skipEmptyLines(record []string) bool {
	return len(record) == 1 && strings.TrimSpace(record[0]) == ""
}

// parseCategoryOrFundHouse processes a record to determine whether it represents
// a scheme category or a fund house. It updates the currentCategory and currentFundHouse accordingly.
func parseCategoryOrFundHouse(record []string, categories *[]SchemeCategory, currentCategory **SchemeCategory, currentFundHouse **FundHouse) {
	if len(record) == 1 {
		// If the record contains "Schemes", it's a new scheme category
		if strings.Contains(record[0], "Schemes") {
			// Finalize the current category and fund house before creating a new category
			finalizeCategoryAndFundHouse(categories, *currentCategory, *currentFundHouse)
			*currentCategory = &SchemeCategory{Name: strings.TrimSpace(record[0])}
			*currentFundHouse = nil
		} else {
			// It's a new fund house, finalize the previous one
			finalizeFundHouse(*currentCategory, *currentFundHouse)
			*currentFundHouse = &FundHouse{Name: strings.TrimSpace(record[0])}
		}
	}
}

// parseFund processes a record representing a mutual fund and appends it to the current fund house.
func parseFund(record []string, currentFundHouse *FundHouse) {
	if currentFundHouse != nil {
		fund := MutualFund{
			SchemeCode:          record[0],
			SchemeName:          record[1],
			ISINDivPayout:       record[2],
			ISINDivReinvestment: record[3],
			NetAssetValue:       record[4],
			RepurchasePrice:     record[5],
			SalePrice:           record[6],
			Date:                record[7],
		}
		currentFundHouse.Funds = append(currentFundHouse.Funds, fund)
	}
}

// finalizeCategoryAndFundHouse ensures the current fund house is added to the current category
// and that the current category is added to the list of categories.
func finalizeCategoryAndFundHouse(categories *[]SchemeCategory, currentCategory *SchemeCategory, currentFundHouse *FundHouse) {
	// Finalize the current fund house if not already done
	finalizeFundHouse(currentCategory, currentFundHouse)
	// Add the category to the list if it has fund houses
	if currentCategory != nil && len(currentCategory.FundHouses) > 0 {
		*categories = append(*categories, *currentCategory)
	}
}

// finalizeFundHouse adds the current fund house to the current category's FundHouses list,
// provided the fund house contains funds.
func finalizeFundHouse(currentCategory *SchemeCategory, currentFundHouse *FundHouse) {
	if currentFundHouse != nil && len(currentFundHouse.Funds) > 0 {
		if currentCategory != nil {
			currentCategory.FundHouses = append(currentCategory.FundHouses, *currentFundHouse)
		}
	}
}

// displayData prints the parsed categories, fund houses, and mutual fund details in a table format.
func displayData(categories []SchemeCategory) {
	for _, category := range categories {
		fmt.Printf("\nScheme Category: %s\n\n", category.Name)

		for _, fundHouse := range category.FundHouses {
			fmt.Printf("Fund House: %s\n", fundHouse.Name)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Scheme Code", "Scheme Name", "ISIN Div Payout", "ISIN Div Reinvestment", "Net Asset Value", "Date"})

			for _, fund := range fundHouse.Funds {
				table.Append([]string{
					fund.SchemeCode,
					fund.SchemeName,
					fund.ISINDivPayout,
					fund.ISINDivReinvestment,
					fund.NetAssetValue,
					fund.Date,
				})
			}

			table.Render()
			fmt.Println() // Add a blank line after each table
		}
	}
}
