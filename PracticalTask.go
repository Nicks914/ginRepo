package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func init() {
	// Initialize MySQL database connection
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/cetec") // "db_username:db_password@tcp(localhost:3306)/cetec"
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/person/:person_id/info", func(c *gin.Context) {
		personID := c.Param("person_id")
		pid, err := strconv.Atoi(personID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
			return
		}

		// Query to fetch person details, phone number, and address
		query := `
			SELECT p.name, ph.number AS phone_number, a.city, a.state, a.street1, a.street2, a.zip_code
			FROM person p
			LEFT JOIN phone ph ON p.id = ph.person_id
			LEFT JOIN address_join aj ON p.id = aj.person_id
			LEFT JOIN address a ON aj.address_id = a.id
			WHERE p.id = ?
		`

		var (
			name    string
			phone   sql.NullString
			city    sql.NullString
			state   sql.NullString
			street1 sql.NullString
			street2 sql.NullString
			zipCode sql.NullString
		)

		err = db.QueryRow(query, pid).Scan(&name, &phone, &city, &state, &street1, &street2, &zipCode)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch person info"})
			return
		}

		// Prepare JSON response
		response := gin.H{
			"name":         name,
			"phone_number": phone.String,
			"city":         city.String,
			"state":        state.String,
			"street1":      street1.String,
			"street2":      street2.String,
			"zip_code":     zipCode.String,
		}

		c.JSON(http.StatusOK, response)
	})

	// Task 2: REST Endpoint (POST) to create a new person
	r.POST("/person/create", func(c *gin.Context) {
		var person struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			City        string `json:"city"`
			State       string `json:"state"`
			Street1     string `json:"street1"`
			Street2     string `json:"street2"`
			ZipCode     string `json:"zip_code"`
		}

		if err := c.BindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert into person table
		result, err := db.Exec(`
			INSERT INTO person(name) VALUES(?)
		`, person.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
			return
		}

		personID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person ID"})
			return
		}

		// Insert into phone table
		_, err = db.Exec(`
			INSERT INTO phone(person_id, number) VALUES(?, ?)
		`, personID, person.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create phone number"})
			return
		}

		// Insert into address table
		result, err = db.Exec(`
			INSERT INTO address(city, state, street1, street2, zip_code) VALUES(?, ?, ?, ?, ?)
		`, person.City, person.State, person.Street1, person.Street2, person.ZipCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
			return
		}

		addressID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve address ID"})
			return
		}

		// Insert into address_join table
		_, err = db.Exec(`
			INSERT INTO address_join(person_id, address_id) VALUES(?, ?)
		`, personID, addressID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link address to person"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
	})

	// Run the server
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
