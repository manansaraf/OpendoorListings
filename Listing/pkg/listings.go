package Listings

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

// This struct contains the filters to be set on the data when it is called.
type ListingFilter struct {
	minPrice int
	maxPrice int
	minBath  int
	maxBath  int
	minBed   int
	maxBed   int
}

/**
*All the below structs are used to create the json return value
 */
type ListingCollection struct {
	Collection string    `json:"type"`
	Listings   []Listing `json:"features"`
}

type Listing struct {
	Type       string            `json:"type"`
	Geometry   Location          `json:"geometry"`
	Properties ListingProperties `json:"properties"`
}

type Location struct {
	LocationType string    `json:"type"`
	Coordinates  []float64 `json:"coordinates"`
}

type ListingProperties struct {
	ID          string `json:"id"`
	Price       int    `json:"price"`
	Address     string `json:"street"`
	NumBedroom  int    `json:"bedrooms"`
	NumBathroom int    `json:"bathrooms"`
	Area        int    `json:"sq_ft"`
}

// This function creates a new Listing filter with default values
func NewListingFilter() ListingFilter {
	filter := ListingFilter{minPrice: 0, maxPrice: math.MaxInt64,
		minBath: 0, maxBath: math.MaxInt64, minBed: 0, maxBed: math.MaxInt64}
	return filter
}

// Resets the filter values after a request has been served
func (filter *ListingFilter) ResetFilter() {
	filter.minPrice = 0
	filter.maxPrice = math.MaxInt64
	filter.minBath = 0
	filter.maxBath = math.MaxInt64
	filter.minBed = 0
	filter.maxBed = math.MaxInt64
}

// This function is the function handler for http requests
func (filter *ListingFilter) ReceiveAndRespondRequest(res http.ResponseWriter, req *http.Request) {
	// Get all the parameter values
	minPrice := req.URL.Query().Get("min_price")
	if len(minPrice) != 0 {
		price, err := strconv.Atoi(minPrice)
		if err != nil {
			fmt.Errorf("Error in reading min_price: %v", minPrice)
		} else {
			filter.minPrice = price
		}
	}
	maxPrice := req.URL.Query().Get("max_price")
	if len(maxPrice) != 0 {
		price, err := strconv.Atoi(maxPrice)
		if err != nil {
			fmt.Errorf("Error in reading max_price: %v", maxPrice)
		} else {
			filter.maxPrice = price
		}
	}
	minBath := req.URL.Query().Get("min_bath")
	if len(minBath) != 0 {
		bath, err := strconv.Atoi(minBath)
		if err != nil {
			fmt.Errorf("Error in reading min_bath: %v", minBath)
		} else {
			filter.minBath = bath
		}
	}
	maxBath := req.URL.Query().Get("max_bath")
	if len(maxBath) != 0 {
		bath, err := strconv.Atoi(maxBath)
		if err != nil {
			fmt.Errorf("Error in reading max_bath: %v", maxBath)
		} else {
			filter.maxBath = bath
		}
	}
	minBed := req.URL.Query().Get("min_bed")
	if len(minBed) != 0 {
		bed, err := strconv.Atoi(minBed)
		if err != nil {
			fmt.Errorf("Error in reading min_bed: %v", minBed)
		} else {
			filter.minBed = bed
		}
	}
	maxBed := req.URL.Query().Get("max_bed")
	if len(maxBed) != 0 {
		bed, err := strconv.Atoi(maxBed)
		if err != nil {
			fmt.Errorf("Error in reading max_bed: %v", maxBed)
		} else {
			filter.maxBed = bed
		}
	}
	// Build the JSON from filter parameters given
	data := filter.buildJSONFromFilter()
	_, err := res.Write(data)
	if err != nil {
		fmt.Errorf("Error in sending data: %v", err)
	}
	filter.ResetFilter()
}

// Given the filter parameter then creates a JSON file to return back to the client
func (filter *ListingFilter) buildJSONFromFilter() []byte {
	var listings []Listing
	// Open csv file
	csvfile, err := os.Open("pkg/listings.csv")
	if err != nil {

		fmt.Errorf("Error in opening csv file: %v", err)
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	counter := 0
	for {
		record, err := reader.Read()
		counter++
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Errorf("Error in reading csv file: %v", err)
			continue
		}
		// Create a listing struct from each record
		listing := filter.createListingFromRecord(record)
		// Check whether the listing comes under the filters
		if filter.checkRecordMatchesFilter(listing) {
			listings = append(listings, listing)
		}
	}
	collection := ListingCollection{Collection: "FeatureCollection", Listings: listings}
	Body, err := json.Marshal(collection)
	if err != nil {
		fmt.Errorf("Error in marshalling json: %v", err)
	}
	return Body
}

// Thus function checks whehter the listing satisfies the parameters given in the GET request
func (filter *ListingFilter) checkRecordMatchesFilter(listing Listing) bool {
	if listing.Properties.Price >= filter.minPrice && listing.Properties.Price <= filter.maxPrice {
		if listing.Properties.NumBedroom >= filter.minBed && listing.Properties.NumBedroom <= filter.maxBed {
			if listing.Properties.NumBathroom >= filter.minBath && listing.Properties.NumBathroom <= filter.maxBath {
				return true
			}
		}
	}
	return false
}

// This function basically read the csv file and creates the listing struct
func (filter *ListingFilter) createListingFromRecord(record []string) Listing {

	price, err := strconv.Atoi(record[3])
	if err != nil {
		fmt.Errorf("Error parsing price: %v", err)
	}
	bed, err := strconv.Atoi(record[4])
	if err != nil {
		fmt.Errorf("Error parsing number of bedrooms: %v", err)
	}
	bath, err := strconv.Atoi(record[5])
	if err != nil {
		fmt.Errorf("Error parsing number of bathrooms: %v", err)
	}
	area, err := strconv.Atoi(record[6])
	if err != nil {
		fmt.Errorf("Error parsing area: %v", err)
	}
	lat, err := strconv.ParseFloat(record[7], 64)
	if err != nil {
		fmt.Errorf("Error parsing latitude: %v", err)
	}
	long, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		fmt.Errorf("Error parsing longitude: %v", err)
	}
	coordinates := []float64{lat, long}
	loc := Location{LocationType: "Point", Coordinates: coordinates}
	properties := ListingProperties{ID: record[0], Price: price, Address: record[1],
		NumBedroom: bed, NumBathroom: bath, Area: area}
	listing := Listing{Type: "Feature", Geometry: loc, Properties: properties}
	return listing
}
