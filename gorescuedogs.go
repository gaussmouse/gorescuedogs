/*
Gorescuedogs fetches and displays dog listing data from Oregon Dog Rescue
using the Petfinder API.

It has three actions that can be accessed using flags (see below).
If no flags are set or an invalid flag is included, gorescuedogs
prints the usage options.

Usage:

	gorescuedogs [options]

Flag options:

	-today
		Fetch dogs posted today.

	-3days
		Fetch dogs posted in the last 3 days.

	-filter
		Filter dogs based on options (age, size, gender).
		When prompted, pick none to all values for each option
		in the listed order (comma-delimited).
		1. age values		= 	baby, young, adult, senior
		2. size values		= 	small, medium, large, xlarge
		3. gender values 	= 	male, female

gorescuedogs does not automatically save the results. In order to save listings
for further reference, include `> outputFileName.txt` after all flags.
If using the filter option, the prompts print to the text file instead of the
console window, but input is still accepted in the same order.
*/
package main

import (
	"fmt"
	"gorescuedogs/cli"
	"gorescuedogs/internals"
)

func main() {
	// Set Petfinder API Key and Secret
	apiKey, apiSecret := internals.GetAPIKeyAndSecret()

	accessToken, err := internals.Authenticate(apiKey, apiSecret)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}

	cli.RunAction(accessToken)
}
