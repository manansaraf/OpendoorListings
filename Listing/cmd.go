package main

import (
	"github.com/OpendoorListings/Listing/pkg"
	"net/http"
)

func main() {
	listingFilter := Listings.NewListingFilter()
	http.HandleFunc("/listings", listingFilter.ReceiveAndRespondRequest)
	http.ListenAndServe(":8000", nil)
}
