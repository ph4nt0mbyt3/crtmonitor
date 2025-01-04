# CrtMonitor

CrtMonitor is a lightweight and efficient Go-based tool for monitoring subdomains of a specific domain using data from crt.sh. It operates without relying on a database, making it simple to set up and use. By leveraging data from certificate transparency logs, it allows users to detect newly issued certificates containing subdomains and compare them against a previously saved list, perfect for security professionals and domain owners who want to keep track of changes in their subdomain infrastructure.

# Key Features:

No Database Dependency: Saves data in lightweight JSON files, eliminating the need for complex setups.
Fetch Subdomains: Retrieves subdomains of a domain from crt.sh with ease.
Monitor Mode: Detects and highlights newly discovered subdomains compared to prior scans.
User-Friendly Output: Displays results in a clean, sorted format.
Portable and Fast: Designed to be lightweight and efficient, ideal for quick checks.

Installation Instructions:
To install crtmonitor, ensure you have Go installed. Then, run the following command:

```
go install github.com/ph4nt0mbyt3/crtmonitor@latest
```
This will download, compile, and place the crtmonitor binary in your $GOPATH/bin directory (usually ~/go/bin). Make sure this directory is added to your PATH environment variable to run the tool from any location.

# Usage:

1- Help

```
crtmonitor -h
Usage of ./crtmonitor:
  -domain string
        Domain identity to check
  -monitor
        Enable monitoring mode
```

2- Basic Usage:
Fetch and display all subdomains of a given domain:

```
crtmonitor -domain example.com
```

3-Monitoring Mode:
Compare current subdomains with a saved list and detect new ones:

```
crtmonitor -domain example.com -monitor
```

While crtmonitor it relies on crt.sh, a third-party service. As a result:
Occasional errors or downtime from crt.sh could affect the tool's functionality.
