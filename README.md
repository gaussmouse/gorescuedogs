# gorescuedogs
gorescuedogs is a Go-based web scraper that retrieves dog listings from Oregon Dog Rescue's website using the Petfinder API.

## Getting Started

To get started with the project, follow these steps:

1. **Get a Petfinder API Key and Secret:**
    To use the Petfinder API, you need to obtain an API key and secret. Visit the [Petfinder Developer Portal](https://www.petfinder.com/developers/) to sign up for an account and create an application. Once registered, you will receive an API key and secret. Keep these credentials secure.

2. **Clone the Repository:**
    Clone this repository to your local machine:
        `git clone https://github.com/yourusername/gorescuedogs.git`
        `cd gorescuedogs`
3. **Add API Credentials to the Code:**
    Open the internals.go file in the internals package and replace the placeholder values in the GetAPIKeyAndSecret function with your actual API key and secret.

4. **Build the Executable:**
    Build the executable using the following command:
        `go build gorescuedogs`

5. **Run the Program:**
    Run the program using the following command:
        `./gorescuedogs [options]`
    By default, the program will output the program usage. 
    Use flags to specify different actions:
    - `-today`: Fetch dogs posted today
    - `-3days`: Fetch dogs posted in the last 3 days
    - `-filter`: Filter dogs based on options

6. **Run Tests:**
    Run unit tests for the cli and internals packages with the following commands:
        `go test ./cli`
        `go test ./internals`

## Binary Download
You can also download the pre-built binary executable from the Releases section. Make sure to give the executable execute permissions before running it:
    `chmod +x gorescuedogs`
    `./gorescuedogs [options]`

## Contributing
Contributions are welcome! If you find any issues or have suggestions, please open an issue or create a pull request.

## License
This project is licensed under the MIT License.