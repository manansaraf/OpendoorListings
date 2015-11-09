package main

import (
	"github.com/OpendoorListings/Listing/pkg"
	"net/http"
	"os"
)

func main() {
	listingFilter := Listings.NewListingFilter()
	http.HandleFunc("/listings", listingFilter.ReceiveAndRespondRequest)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
