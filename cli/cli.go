// Package cli provides command-line interface functionality for the
// gorescuedogs program.
//
// It handles user interactions through flags, processes different actions,
// and interfaces with the Petfinder API.
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"gorescuedogs/internals"
	"os"
	"strings"
)

// PrintUsage prrints the usage options for the program
// Also used to customize output for flag.Usage func()
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gorescuedogs [options]")
	fmt.Println("  -today")
	fmt.Println("\tFetch dogs posted today")
	fmt.Println("  -3days")
	fmt.Println("\tFetch dogs posted in the last 3 days")
	fmt.Println("  -filter")
	fmt.Println("\tFilter dogs based on options")
}

// processAction handles the processing logic for each selected action.
func processAction(apiURL, noDogsMessage string, accessToken string) {
	response, err := internals.FetchAnimals(apiURL, accessToken)
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

// RunAction actives program action based on the flag(s) set by the user
// Prompts PrintUsage() if no flags are passed
func RunAction(accessToken string) {
	// Define command-line flags for different actions
	todayFlag := flag.Bool("today", false, "Fetch dogs posted today")
	threeDaysFlag := flag.Bool("3days", false, "Fetch dogs posted in the last 3 days")
	filterFlag := flag.Bool("filter", false, "Filter dogs based on options")

	// Set custom usage message
	flag.Usage = func() {
		PrintUsage()
	}

	// Parse the command-line flags
	flag.Parse()

	// If none of the valid flags are set, print flag options as reminder
	if !*todayFlag && !*threeDaysFlag && !*filterFlag {
		PrintUsage()
		os.Exit(1)
	}

	if *todayFlag {
		fmt.Println("Looking for new dogs posted today...")
		apiURL := internals.GenerateAPIURL("TODAY")
		processAction(apiURL, "No new dogs today :(", accessToken)
	}
	if *threeDaysFlag {
		fmt.Println("Looking for new dogs posted in the last 3 days...")
		apiURL := internals.GenerateAPIURL("3DAYS")
		processAction(apiURL, "No new dogs posted recently :(", accessToken)
	}
	if *filterFlag {
		fmt.Println("*---------------------------------------*")
		fmt.Println("|            Filter Options             |")
		fmt.Println("| (Enter none to all for each category) |")
		fmt.Println("|                                       |")
		fmt.Println("|  Age:    baby, young, adult, senior   |")
		fmt.Println("|  Size:   small, medium, large, xlarge |")
		fmt.Println("|  Gender: male, female                 |")
		fmt.Println("*---------------------------------------*")
		// Show the filter options and get user inputs for each filter category
		ageOptions := getUserFilterOptions("Age", []string{"baby", "young", "adult", "senior"})
		sizeOptions := getUserFilterOptions("Size", []string{"small", "medium", "large", "xlarge"})
		genderOptions := getUserFilterOptions("Gender", []string{"male", "female"})

		fmt.Println("Looking for new dogs with the selected filter options:", ageOptions, sizeOptions, genderOptions)

		// Construct the filtered API URL based on user-selected options
		filteredAPIURL := internals.GenerateFilteredAPIURL(ageOptions, sizeOptions, genderOptions)

		processAction(filteredAPIURL, "No dogs match the selected filters :(", accessToken)
	}
}

// processData processes the fetched animal data and prints relevant information.
func processData(animals []internals.Dog) {
	fmt.Println("-----------------------")
	for _, animal := range animals {
		fmt.Println("Name:", animal.Name)
		fmt.Println("Age:", animal.Age)
		fmt.Println("Gender:", animal.Gender)
		fmt.Println("Size:", animal.Size)
		fmt.Println("Breed:", animal.Breeds.Primary)
		fmt.Println("URL:", animal.URL)
		fmt.Println("-----------------------")
	}
}

// getUserFilterOptions prompts the user to select options for a specific filter category.
// The user input is validated against the provided options, and the selected options are returned as a comma-separated string.
func getUserFilterOptions(category string, options []string) string {
	userInput := GetUserInput(fmt.Sprintf("Enter %s options: ", category))

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

// GetUserInput prompts the user with the specified message and reads the input from the command-line.
func GetUserInput(promptMessage string) string {
	fmt.Print(promptMessage)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// Trim any leading/trailing spaces and newline characters from the input
	input = strings.TrimSpace(input)

	return input
}
