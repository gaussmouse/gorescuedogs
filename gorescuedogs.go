package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Struct to hold the response from the Petfinder API token request.
type tokenResponse struct {
	Token string `json:"access_token"`
}

// Struct to represent a dog from the Petfinder API response.
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
	apiKey, apiSecret := getAPIKeyAndSecret()

	accessToken, err := authenticate(apiKey, apiSecret)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	runActionLoop(accessToken)
}

// processAction handles the processing logic for each selected action.
func processAction(action, apiURL, noDogsMessage string, accessToken string) {
	response, err := fetchAnimals(apiURL, accessToken)
	if err != nil {
		fmt.Println("Failed to fetch animals:", err)
		return
	}
	if len(response.Animals) == 0 {
		fmt.Println(noDogsMessage)
	} else {
		processData(response.Animals)
	}
}

// runActionLoop is responsible for repeatedly presenting action options to the user,
// processing the user's choice, and interacting with the Petfinder API based on the chosen action.
// The loop continues until the user chooses to exit by pressing 'q' or 'Q'.
func runActionLoop(accessToken string) {
	for {
		// Get the action from the user
		action := getUserInput("Select an action (TODAY/3DAYS/FILTER), or press Q to quit: ")

		if strings.EqualFold(action, "q") {
			fmt.Println("Exiting the program.")
			break
		}

		switch strings.ToUpper(action) {
		case "TODAY":
			apiURL := generateAPIURL("TODAY")
			processAction(action, apiURL, "No new dogs posted today.", accessToken)

		case "3DAYS":
			apiURL := generateAPIURL("3DAYS")
			processAction(action, apiURL, "No new dogs posted in the last 3 days.", accessToken)

		case "FILTER":
			// Show the filter options and get user inputs for each filter category
			ageOptions := getUserFilterOptions("Age", []string{"baby", "young", "adult", "senior"})
			sizeOptions := getUserFilterOptions("Size", []string{"small", "medium", "large", "xlarge"})
			genderOptions := getUserFilterOptions("Gender", []string{"male", "female"})

			// Construct the filtered API URL based on user-selected options
			filteredAPIURL := generateFilteredAPIURL(ageOptions, sizeOptions, genderOptions)

			// Process the filtered action
			processAction(action, filteredAPIURL, "No dogs match the selected filters.", accessToken)

		default:
			fmt.Println("Invalid action. Please select one of the provided actions.")
		}
	}
}

// getUserFilterOptions prompts the user to select options for a specific filter category.
// The user input is validated against the provided options, and the selected options are returned as a comma-separated string.
func getUserFilterOptions(category string, options []string) string {
	fmt.Printf("Select %s options (comma-separated): %s\n", category, strings.Join(options, ", "))
	userInput := getUserInput(fmt.Sprintf("Enter %s options: ", category))

	// Split the user input by commas
	selectedOptions := strings.Split(userInput, ",")

	// Trim leading/trailing spaces from each selected option
	for i := range selectedOptions {
		selectedOptions[i] = strings.TrimSpace(selectedOptions[i])
	}

	// Filter out any invalid options
	validOptions := make([]string, 0)
	for _, option := range selectedOptions {
		for _, validOption := range options {
			if strings.EqualFold(option, validOption) {
				validOptions = append(validOptions, option)
				break
			}
		}
	}

	// Return the selected options as a comma-separated string
	return strings.Join(validOptions, ",")
}

// getAPIKeyAndSecret retrieves the API key and secret.
func getAPIKeyAndSecret() (string, string) {
	// Replace these values with your actual API key and secret.
	apiKey := "API Key goes here"
	apiSecret := "API Secret goes here"
	return apiKey, apiSecret
}

// authenticate requests an access token from the Petfinder API using the provided API key and secret.
// It returns the access token and an error if the authentication fails.
func authenticate(apiKey, apiSecret string) (string, error) {
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
	var tokenResp tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}

// generateAPIURL constructs the Petfinder API URL with the current date as the "after" parameter.
func generateAPIURL(dateRange string) string {
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

// generateFilteredAPIURL constructs the Petfinder API URL with the current date as the "after" parameter
// and filters based on the user-selected options.
func generateFilteredAPIURL(ageOptions, sizeOptions, genderOptions string) string {
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

// fetchAnimals makes a request to the Petfinder API to fetch adoptable dogs.
// It takes the API URL and access token as input and returns a pointer to the struct
// containing the response data and an error (if any).
func fetchAnimals(apiURL, accessToken string) (*struct {
	Animals []dog `json:"animals"`
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
		Animals []dog `json:"animals"`
	}
	err = json.NewDecoder(apiResp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Return the pointer to the struct containing the fetched animal data.
	return &response, nil
}

// processData processes the fetched animal data and prints relevant information.
func processData(animals []dog) {
	fmt.Println("-----------------------")
	for _, animal := range animals {
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

// getUserInput prompts the user with the specified message and reads the input from the command-line.
func getUserInput(promptMessage string) string {
	fmt.Print(promptMessage)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// Trim any leading/trailing spaces and newline characters from the input
	input = strings.TrimSpace(input)

	return input
}
