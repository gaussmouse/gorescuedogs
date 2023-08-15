package internals

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// Test the GenerateAPIURL function with the "TODAY" date range.
// Check if the generated URL contains the correct query parameters
// for today's listings.
func TestGenerateAPIURLToday(t *testing.T) {
	// Get the current date
	now := time.Now()

	// Get the beginning of the day (00:00:00) for the current date
	beginningOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Format the beginning of the day in ISO8601 format
	beginningOfDayStr := beginningOfDay.Format("2006-01-02T15:04:05-07:00")

	// Construct the expected API URL
	want := "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&after=" + url.QueryEscape(beginningOfDayStr)

	// Define the time duration for calculating "after" parameter value
	timeDuration := "TODAY"

	// Call the function
	got := GenerateAPIURL(timeDuration)

	// Assert that the generated API URL matches the expected value
	if got != want {
		t.Errorf("TestGenerateAPIURLToday failed.\nGot:  %s\nWant: %s", got, want)
	}
}

// Test the GenerateAPIURL function with the "3DAYS" date range.
// Check if the generated URL contains the correct query parameters for listings in the last 3 days.
func TestGenerateAPIURL3Days(t *testing.T) {
	// Get the current date
	now := time.Now()

	// Calculate the date 2 days ago
	twoDaysAgo := now.AddDate(0, 0, -2)

	// Get the beginning of the day (00:00:00) for the calculated date
	beginningOfDay := time.Date(twoDaysAgo.Year(), twoDaysAgo.Month(), twoDaysAgo.Day(), 0, 0, 0, 0, twoDaysAgo.Location())

	// Format the beginning of the day in ISO8601 format
	beginningOfDayStr := beginningOfDay.Format("2006-01-02T15:04:05-07:00")

	// Construct the expected API URL
	want := "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&after=" + url.QueryEscape(beginningOfDayStr)

	// Define the time duration for calculating "after" parameter value
	timeDuration := "3DAYS"

	// Call the function
	got := GenerateAPIURL(timeDuration)

	// Assert that the generated API URL matches the expected value
	if got != want {
		t.Errorf("TestGenerateAPIURL3Days failed.\nGot:  %s\nWant: %s", got, want)
	}
}

// Test the GenerateFilteredAPIURL function with various filter options.
// Check if the generated URL includes the selected filter options in the query parameters.
func TestGenerateFilteredAPIURL(t *testing.T) {
	testCases := []struct {
		ageOptions    string
		sizeOptions   string
		genderOptions string
		want          string
	}{
		// Test case 1: All options selected
		{
			ageOptions:    "baby,young",
			sizeOptions:   "small,medium",
			genderOptions: "male,female",
			want:          "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&age=baby%2Cyoung&size=small%2Cmedium&gender=male%2Cfemale",
		},
		// Test case 2: Only age options selected
		{
			ageOptions:    "adult,senior",
			sizeOptions:   "",
			genderOptions: "",
			want:          "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&age=adult%2Csenior",
		},
		// Test case 3: Only size options selected
		{
			ageOptions:    "",
			sizeOptions:   "large,xlarge",
			genderOptions: "",
			want:          "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&size=large%2Cxlarge",
		},
		// Test case 4: Only gender options selected
		{
			ageOptions:    "",
			sizeOptions:   "",
			genderOptions: "male",
			want:          "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable&gender=male",
		},
		// Test case 5: No options selected
		{
			ageOptions:    "",
			sizeOptions:   "",
			genderOptions: "",
			want:          "https://api.petfinder.com/v2/animals?type=dog&organization=OR208&status=adoptable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.ageOptions+"_"+tc.sizeOptions+"_"+tc.genderOptions, func(t *testing.T) {
			got := GenerateFilteredAPIURL(tc.ageOptions, tc.sizeOptions, tc.genderOptions)
			if got != tc.want {
				t.Errorf("TestGenerateFilteredAPIURL failed.\nGot:  %s\nWant: %s", got, tc.want)
			}
		})
	}
}

// Mock HTTP server for testing with TestFetchAnimals
func mockAPIHandler(responseBody string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, responseBody)
	}
}

// Test the FetchAnimals function by mocking a response from the Petfinder API.
// Verify if the fetched data is correctly processed and returned in the response.
func TestFetchAnimals(t *testing.T) {
	// Create a mock server and set its URL as the API endpoint
	mockServer := httptest.NewServer(mockAPIHandler(`{"animals": [{"id": 1, "name": "Dog"}]}`))
	defer mockServer.Close()

	accessToken := "mockAccessToken"

	got, err := FetchAnimals(mockServer.URL, accessToken)
	if err != nil {
		t.Fatalf("fetchAnimals failed with error: %v", err)
	}

	want := &struct {
		Animals []Dog `json:"animals"`
	}{
		Animals: []Dog{{ID: 1, Name: "Dog"}},
	}

	if len(got.Animals) != len(want.Animals) {
		t.Errorf("TestFetchAnimals failed.\nGot:  %v\nWant: %v", got, want)
	}
}
