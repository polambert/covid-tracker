
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"os"

	"github.com/gin-gonic/gin"
)

//
type CountryResponse struct {
	Country string 				`json:"country"`
	Cases uint32 				`json:"cases"`
	TodayCases uint32 			`json:"todayCases"`
	Deaths uint32				`json:"deaths"`
	TodayDeaths uint32			`json:"todayDeaths"`
	Recovered uint32			`json:"recovered"`
	Active uint32				`json:"active"`
	Critical uint32				`json:"critical"`
	CasesPerOneMillion uint32 	`json:"casesPerOneMillion"`
	DeathsPerOneMillion uint32	`json:"deathsPerOneMillion"`
	TotalTests uint32			`json:"totalTests"`
	TestsPerOneMillion uint32	`json:"testsPerOneMillion"`
}

//
const apiPath = "https://coronavirus-19-api.herokuapp.com"

var currentCountryData []CountryResponse

//
func check(err error) {
	if (err != nil) {
		fmt.Println("ERR: " + err.Error())
	}
}

func dataAgeByMinutes() float64 {
	file, err := os.Stat("data.json")
	check(err)

	return time.Now().Sub(file.ModTime()).Minutes()
}

func fetchData() {
	// First, check if the data is old enough to need to be updated
	if (dataAgeByMinutes() < 120) {
		// Data isn't old enough
		return
	}

	// Data is old enough
	// stores data into a file
	resp, err := http.Get(apiPath + "/countries/")
	check(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	// turn response data into a readable interface
	// body is already a []byte, which json.Unmarshal requires
	err = json.Unmarshal(body, &currentCountryData)

	if (err != nil) {
		// stop execution here so useless data isn't written to file
		fmt.Println("ERROR: Unable to unmarshal JSON response correctly")
		return
	}

	// Assume data has been extracted correctly, so save to file
	// We're circling back and using the raw byte slice 'body' instead of
	//  wasting time converting currentCountryData back to bytes to write
	//  to file.
	f, err := os.Create("data.json")
	check(err)

	defer f.Close()

	nBytesWritten, err := f.Write(body)
	check(err)

	f.Sync()

	fmt.Printf("Wrote %d bytes to file\n", nBytesWritten)
}

func fetchDataFromFile() int {
	data, err := ioutil.ReadFile("data.json")
	check(err)

	err = json.Unmarshal(data, &currentCountryData)
	if (err != nil) {
		fmt.Println("ERROR: Unable to load country data from disk. Exiting")
		return -1
	}

	return 0 // success
}

func main() {
	fmt.Println("Starting...")

	// backend setup
	errorCode := fetchDataFromFile()

	if (errorCode != 0) {
		fmt.Println("Got nonzero from fetchDataFromFile(), exiting")
		return
	}

	fmt.Println("Grabbed data successfully")

	// web setup
	router := gin.Default()

	router.Static("/assets", "./assets") // give public access to assets folder
	router.LoadHTMLGlob("./templates/*")

	router.GET("/", func(c *gin.Context) {
		fetchData()

		c.HTML(200, "index.html", gin.H{
			"Data": currentCountryData,
			"MinutesAgo": int(dataAgeByMinutes()),
		})
	})

	router.Run(":8080")
}
