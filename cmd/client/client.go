package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/dominik-matic/dddns-apiclient/internal/apiclient"
)

func main() {
	tokenPath := flag.String("token", apiclient.DEFAULT_TOKEN_PATH, "Path to the auth token file")
	mode := flag.String("mode", apiclient.DEFAULT_MODE, "Whether to update or delete the record. Options: "+strings.Join(apiclient.ALLOWED_MODES, ", "))

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: missing required positional argument")
		fmt.Println("Usage: client [options] <space-separated-domain-names>")
		fmt.Println("Examples:")
		fmt.Printf("\tclient mysubdomain.mydomain.com othersubdomain.mydomain.com\n")
		fmt.Printf("\tclient -token=~/.my_token -mode=delete mysubdomain.mydomain.com othersubdomain.mydomain.com\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	token, err := loadToken(*tokenPath)
	if err != nil {
		fmt.Printf("Token error: %v\n", err)
		os.Exit(2)
	}
	if ok := validateMode(*mode); !ok {
		fmt.Printf("Unknown mode: %v\n", *mode)
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(3)
	}

	// TODO: update apiserver to be able to handle a list a domains
	// only take into account the 1st domain for now
	domain := args[0]
	if err = apiclient.SendRequest(token, *mode, domain); err != nil {
		fmt.Printf("Request failed: %v", err)
		os.Exit(4)
	}
	fmt.Println("Request Sent Successfully :D")
}

func loadToken(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, nil
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", fmt.Errorf("no valid auth token found in file")
}

func validateMode(mode string) bool {
	return slices.Contains(apiclient.ALLOWED_MODES, mode)
}
