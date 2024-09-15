# mf-nav-parser

> [!NOTE]  
> This is a hobby/learning project, created to explore Go, HTTP requests, CSV parsing, and table formatting. The code is not intended for production use, and the AMFI India website might have restrictions or rate limits on scraping, so please use it responsibly.

## Overview

`mf-nav-parser` is a simple Go program that fetches and parses Mutual Fund Net Asset Value (NAV) data from the AMFI India portal. It organizes the data by scheme categories, fund houses, and mutual fund schemes, and then displays the parsed data in a well-formatted table.

The code retrieves the NAV data in CSV format, parses it, and groups the funds based on their scheme category and fund house. The results are printed in a table format using the [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) package.

## Features

- Fetches the NAV data for mutual funds from the AMFI India website.
- Parses the CSV data into categories, fund houses, and individual mutual funds.
- Displays the data in a table format for better readability.
- Organizes data by:
  - Scheme Category
  - Fund House
  - Mutual Fund schemes (including NAV, ISIN codes, repurchase, and sale prices).

## Requirements

- Go 1.23+ 
- Internet connection to fetch NAV data from the AMFI India portal.

## Dependencies

The project uses the following Go package:

- [`olekukonko/tablewriter`](https://github.com/olekukonko/tablewriter): For rendering tables in the terminal.

You can install the package using `go get`:

```bash
go get github.com/olekukonko/tablewriter
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/mf-nav-parser.git
cd mf-nav-parser
```

2. Install the dependencies:

```bash
go get ./...
```

3. Run the program:

```bash
go run main.go
```

## How It Works

1. **Fetching Data:**  
   The program fetches the latest NAV data from the AMFI India portal. By default, it retrieves the data for the previous day.
   
2. **Parsing the CSV:**  
   The CSV file is processed line by line, and the following structure is used to organize the data:
   
   - `SchemeCategory`: A broad category under which mutual funds fall.
   - `FundHouse`: Represents a specific fund house under a scheme category.
   - `MutualFund`: Contains details of each fund like Scheme Code, Scheme Name, ISIN codes, NAV, and date.

3. **Displaying Data:**  
   After parsing, the data is displayed using tables categorized by Scheme Category and Fund House.

## Example Output

After running the program, you will see an output similar to:

```
Scheme Category: Schemes - Equity

Fund House: ABC Mutual Fund
+-------------+-------------------+-----------------+---------------------+----------------+------------+
| SCHEME CODE |    SCHEME NAME     | ISIN DIV PAYOUT | ISIN DIV REINVESTMENT| NET ASSET VALUE|    DATE    |
+-------------+-------------------+-----------------+---------------------+----------------+------------+
| 12345       | ABC Equity Fund    | INE123456789    | INE1234567890        | 123.45         | 14-Sep-2023|
| 67890       | ABC Bluechip Fund  | INE098765432    | INE0987654321        | 98.76          | 14-Sep-2023|
+-------------+-------------------+-----------------+---------------------+----------------+------------+

Fund House: XYZ Mutual Fund
+-------------+--------------------+-----------------+---------------------+----------------+------------+
| SCHEME CODE |    SCHEME NAME      | ISIN DIV PAYOUT | ISIN DIV REINVESTMENT| NET ASSET VALUE|    DATE    |
+-------------+--------------------+-----------------+---------------------+----------------+------------+
| 54321       | XYZ Growth Fund     | INE543216789    | INE5432167890        | 543.21         | 14-Sep-2023|
| 98765       | XYZ Value Fund      | INE678965432    | INE6789654321        | 432.10         | 14-Sep-2023|
+-------------+--------------------+-----------------+---------------------+----------------+------------+
```

## Project Structure

```
.
├── main.go           # Main program logic
└── go.mod            # Go module file
```

## Disclaimer

This project is for educational purposes only. The AMFI India website might have restrictions or rate limits, so please use this tool responsibly. The data format or the URL may change, which could break the functionality of this program.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

