package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

func getSubdomains(identity string) (map[string]bool, error) {
	url := fmt.Sprintf("https://crt.sh/?Identity=%s&output=json", identity)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to retrieve data: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	baseDomain := identity
	subdomainPattern := regexp.MustCompile(fmt.Sprintf(`\b([a-zA-Z0-9-]+\.)+%s\b`, regexp.QuoteMeta(baseDomain)))

	subdomains := make(map[string]bool)
	for _, entry := range data {
		domain, ok := entry["name_value"].(string)
		if !ok {
			continue
		}

		matches := subdomainPattern.FindAllString(domain, -1)
		for _, match := range matches {
			subdomains[match] = true
		}
	}

	return subdomains, nil
}

func saveSubdomains(filename string, subdomains map[string]bool) error {
	list := make([]string, 0, len(subdomains))
	for subdomain := range subdomains {
		list = append(list, subdomain)
	}

	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func loadSubdomains(filename string) (map[string]bool, error) {
	subdomains := make(map[string]bool)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return subdomains, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, subdomain := range list {
		subdomains[subdomain] = true
	}

	return subdomains, nil
}

func main() {
	domain := flag.String("domain", "", "Domain identity to check")
	monitor := flag.Bool("monitor", false, "Enable monitoring mode")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Please provide a domain using the -domain flag.")
		os.Exit(1)
	}

	currentSubdomains, err := getSubdomains(*domain)
	if err != nil {
		fmt.Printf("Error fetching subdomains: %v\n", err)
		os.Exit(1)
	}

	if len(currentSubdomains) == 0 {
		fmt.Println("No subdomains found.")
		return
	}

	if *monitor {
		fileName := strings.ReplaceAll(*domain, ".", "_") + "_subdomains.json"
		oldSubdomains, err := loadSubdomains(fileName)
		if err != nil {
			fmt.Printf("Error loading old subdomains: %v\n", err)
			os.Exit(1)
		}

		newSubdomains := make([]string, 0)
		for subdomain := range currentSubdomains {
			if !oldSubdomains[subdomain] {
				newSubdomains = append(newSubdomains, subdomain)
			}
		}

		if len(newSubdomains) > 0 {
			fmt.Println("New subdomains found:")
			sort.Strings(newSubdomains)
			for _, subdomain := range newSubdomains {
				fmt.Printf("https://%s\n", subdomain)
			}
		} else {
			fmt.Println("No new subdomains found.")
		}

		if err := saveSubdomains(fileName, currentSubdomains); err != nil {
			fmt.Printf("Error saving subdomains: %v\n", err)
		}
	} else {
		fmt.Println("Subdomains found:")
		sortedSubdomains := make([]string, 0, len(currentSubdomains))
		for subdomain := range currentSubdomains {
			sortedSubdomains = append(sortedSubdomains, subdomain)
		}
		sort.Strings(sortedSubdomains)
		for _, subdomain := range sortedSubdomains {
			fmt.Printf("https://%s\n", subdomain)
		}
	}
}
