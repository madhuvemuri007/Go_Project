package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/submit", submitHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Connect to MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/electricvehicle_dataset")
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert data into database
	stmt, err := db.Prepare("INSERT INTO Vehicle (VIN, County, City, State, Postal_Code, Model_Year, Make, Model, Electric_Vehicle_Type, CAFV_Eligibility, Electric_Range, Base_MSRP, Legislative_District, DOL_Vehicle_ID, Vehicle_Location, Electric_Utility, Census_Tract_2020) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		r.FormValue("vin"),
		r.FormValue("county"),
		r.FormValue("city"),
		r.FormValue("state"),
		r.FormValue("postal_code"),
		r.FormValue("model_year"),
		r.FormValue("make"),
		r.FormValue("model"),
		r.FormValue("ev_type"),
		r.FormValue("cafve"),
		r.FormValue("electric_range"),
		r.FormValue("base_msrp"),
		r.FormValue("legislative_district"),
		r.FormValue("dol_vehicle_id"),
		r.FormValue("vehicle_location"),
		r.FormValue("electric_utility"),
		r.FormValue("census_tract"),
	)

	if err != nil {
		fmt.Println("Error executing query:", err.Error()) // Print detailed error message
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data inserted successfully")
}
