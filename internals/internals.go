// Package internals provides internal functionalities for the
// gorescuedogs program.
//
// It includes functions related to API interaction, data structuring,
// and data processing.
//
// A PetFinder API Key and API Secret must be added to
// func GetAPIKeyAndSecret() to successfully run the program.
// To get started, go to https://www.petfinder.com/developers/
// and click the button "GET AN API KEY" under the Start Using
// the API heading.
package internals

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A TokenResponse from the Petfinder API token request.
type TokenResponse struct {
	Token string `json:"access_token"`
}

// A Dog represents a single listing from the Petfinder API response.
// See Petfinder API documentation for all possible response fields.
type Dog struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Age    string `json:"age"`
	Gender string `json:"gender"`
	Size   string `json:"size"`
	Breeds struct {
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
	} `json:"breeds"`
	URL string `json:"url"`
}

// GetAPIKeyAndSecret retrieves the API key and secret.
func GetAPIKeyAndSecret() (string, string) {
	// Replace these values with your actual API key and secret.
	apiKey := "API KEY OR ENV VARIABLE HERE"
	apiSecret := "API SECRET OR ENV VARIABLE HERE"
	return apiKey, apiSecret
}

// Authenticate requests an access token from the Petfinder API using the provided API key and secret.
// It returns the access token and an error if the authentication fails.
func Authenticate(apiKey, apiSecret string) (string, error) {
	authURL := "https://api.petfinder.com/v2/oauth2/token"
	client := &http.Client{}

	// Create the data to be sent in the request body for authentication.
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", apiKey)
	data.Set("client_secret", apiSecret)

	// Create a new HTTP POST request with the authentication URL and the encoded data in the request body.
	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the authentication request to the Petfinder API.
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code to ensure successful authentication.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	// Parse the response body to obtain the access token.
	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}

// GenerateAPIURL constructs the Petfinder API URL with the current date as the "after" parameter.
func GenerateAPIURL(dateRange string) string {
	apiURL := "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable"

	now := time.Now()
	var offset time.Duration

	switch dateRange {
	case "TODAY":
		offset = 0
	case "3DAYS":
		offset = -2 * 24 * time.Hour
	default:
		// Handle invalid date range
		return apiURL
	}

	// Calculate the target date by subtracting the offset from the current time
	targetDate := now.Add(offset)

	// Get the beginning of the day (00:00:00) for the target date
	beginningOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())

	// Format the beginning of the day in ISO8601 format
	beginningOfDayStr := beginningOfDay.Format("2006-01-02T15:04:05-07:00")

	// Append the formatted date and time to the API URL
	apiURL += "&after=" + url.QueryEscape(beginningOfDayStr)

	return apiURL
}

// GenerateFilteredAPIURL constructs the Petfinder API URL with the current date as the "after" parameter
// and filters based on the user-selected options.
func GenerateFilteredAPIURL(ageOptions, sizeOptions, genderOptions string) string {
	apiURL := "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable"

	// Append the user-selected filter options to the API URL if they are not empty
	if ageOptions != "" {
		apiURL += "&age=" + url.QueryEscape(ageOptions)
	}
	if sizeOptions != "" {
		apiURL += "&size=" + url.QueryEscape(sizeOptions)
	}
	if genderOptions != "" {
		apiURL += "&gender=" + url.QueryEscape(genderOptions)
	}

	return apiURL
}

// FetchAnimals makes a request to the Petfinder API to fetch adoptable dogs.
// It takes the API URL and access token as input and returns a pointer to the struct
// containing the response data and an error (if any).
func FetchAnimals(apiURL, accessToken string) (*struct {
	Animals []Dog `json:"animals"`
}, error) {
	client := &http.Client{}

	// Create a new HTTP GET request with the API URL.
	apiReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the access token to authenticate the request.
	apiReq.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the API request to the Petfinder API.
	apiResp, err := client.Do(apiReq)
	if err != nil {
		return nil, err
	}
	defer apiResp.Body.Close()

	// Check the response status code to ensure a successful API request.
	if apiResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", apiResp.StatusCode)
	}

	// Parse the response body into a struct to extract the animal data.
	var response struct {
		Animals []Dog `json:"animals"`
	}
	err = json.NewDecoder(apiResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Return the pointer to the struct containing the fetched animal data.
	return &response, nil
}
