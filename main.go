package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"

	//"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// API endpoint URL to retrieve the dataset test
	//url := "https://data.wa.gov/api/views/f6w7-q2d2/rows.csv?accessType=DOWNLOAD"

	// File path to save the downloaded dataset
	filePath := "C:/Users/madhu/OneDrive/Desktop/Go_Project/dataset.csv"

	// // Create the output file
	// outputFile, err := os.Create(filePath)
	// if err != nil {
	// 	fmt.Printf("Failed to create file: %v\n", err)
	// 	return
	// }
	// defer outputFile.Close()

	// // Make a GET request to the API endpoint
	// response, err := http.Get(url)
	// if err != nil {
	// 	fmt.Printf("Failed to make request: %v\n", err)
	// 	return
	// }
	// defer response.Body.Close()

	// // Check the status code of the response
	// if response.StatusCode != http.StatusOK {
	// 	fmt.Printf("Unexpected status code: %d\n", response.StatusCode)
	// 	return
	// }

	// // Copy the response body to the output file
	// _, err = io.Copy(outputFile, response.Body)
	// if err != nil {
	// 	fmt.Printf("Failed to copy data to file: %v\n", err)
	// 	return
	// }

	// fmt.Println("Dataset downloaded successfully!")

	// Open the downloaded CSV file
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	// Establish a connection to the MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/electricvehicle_dataset")
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	// Count the existing records in the database
	var existingRecords int
	err = db.QueryRow("SELECT COUNT(*) FROM Vehicle").Scan(&existingRecords)
	if err != nil {
		fmt.Printf("Failed to count existing records: %v\n", err)
		return
	}
	fmt.Println(existingRecords)

	// Iterate through each row in the CSV file and insert into the database if it's new
	insertedRecords := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Failed to read CSV record: %v\n", err)
			return
		}

		// Check if the VIN exists in the database
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM Vehicle WHERE VIN = ?", record[0]).Scan(&count)
		if err != nil {
			fmt.Printf("Failed to check existing record: %v\n", err)
			return
		}
		//fmt.Println(count)

		// If the VIN does not exist, insert the record into the database
		if count == 0 {
			_, err = db.Exec("INSERT INTO Vehicle (VIN, County, City, State, Postal_Code, Model_Year, Make, Model, Electric_Vehicle_Type, CAFV_Eligibility, Electric_Range, Base_MSRP, Legislative_District, DOL_Vehicle_ID, Vehicle_Location, Electric_Utility, Census_Tract_2020) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10], record[11], record[12], record[13], record[14], record[15], record[16])
			if err != nil {
				fmt.Printf("Failed to insert data into database: %v\n", err)
				return
			}
			insertedRecords++
		}
	}

	fmt.Printf("%d new records inserted into MySQL database!\n", insertedRecords)
}
