package cli

import (
	"fmt"
	"gorescuedogs/internals"
	"io/ioutil"
	"os"
	"testing"
)

// Test the processData function by providing sample dog data.
// Verify that the function correctly formats and prints the dog data.
func TestProcessData(t *testing.T) {
	// Create a mock dog for testing
	mockDog := internals.Dog{
		Name:   "Test Dog",
		Age:    "Adult",
		Gender: "Male",
		Size:   "Medium",
		Breeds: struct {
			Primary   string `json:"primary"`
			Secondary string `json:"secondary"`
		}{Primary: "Mixed"},
		URL: "https://example.com/dog1",
	}

	// Capture stdout for testing
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
		w.Close()
	}()

	// Call the function being tested with the mock dog
	processData([]internals.Dog{mockDog})

	w.Close()

	// Read the captured output
	capturedOutput, _ := ioutil.ReadAll(r)

	// Define how the output should look like based on the mock dog
	expectedOutput := fmt.Sprintf(
		"-----------------------\nName: %s\nAge: %s\nGender: %s\nSize: %s\nBreed: %s\nURL: %s\n-----------------------\n",
		mockDog.Name, mockDog.Age, mockDog.Gender, mockDog.Size, mockDog.Breeds.Primary, mockDog.URL,
	)

	if string(capturedOutput) != expectedOutput {
		t.Errorf("TestProcessData failed.\nGot:\n%s\nWant:\n%s", string(capturedOutput), expectedOutput)
	}
}
