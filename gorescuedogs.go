package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

type dog struct {
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

func main() {
	// Set Petfinder API Key and Secret
	apiKey := "API Key goes here"
	apiSecret := "API Secret goes here"

	// Set the authentication endpoint URL
	authURL := "https://api.petfinder.com/v2/oauth2/token"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP POST request to obtain the access token
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", apiKey)
	data.Set("client_secret", apiSecret)

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Failed to create authentication request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the authentication request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Authentication failed with status code:", resp.StatusCode)
		return
	}

	// Parse the response body to get the access token
	var tokenResp tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		fmt.Println("Failed to parse authentication response:", err)
		return
	}

	accessToken := tokenResp.Token

	// Make an API request with the access token
	apiURL := "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable"

	apiReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Failed to create API request:", err)
		return
	}

	apiReq.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the API request
	apiResp, err := client.Do(apiReq)
	if err != nil {
		fmt.Println("Failed to make API request:", err)
		return
	}
	defer apiResp.Body.Close()

	// Check the response status code
	if apiResp.StatusCode != http.StatusOK {
		fmt.Println("API request failed with status code:", apiResp.StatusCode)
		return
	}

	// Parse the response body into a struct
	var response struct {
		Animals []dog `json:"animals"`
	}
	err = json.NewDecoder(apiResp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Failed to parse API response body:", err)
		return
	}

	// Process the extracted data
	for _, animal := range response.Animals {
		fmt.Println("ID:", animal.ID)
		fmt.Println("Name:", animal.Name)
		fmt.Println("Age:", animal.Age)
		fmt.Println("Gender:", animal.Gender)
		fmt.Println("Size:", animal.Size)
		fmt.Println("Breed:", animal.Breeds.Primary)
		fmt.Println("URL:", animal.URL)
		fmt.Println("-----------------------")
	}
}
