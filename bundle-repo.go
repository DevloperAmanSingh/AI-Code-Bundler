package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AkhilSharma90/AI-Code-Bundler/internal/files"
	"github.com/AkhilSharma90/AI-Code-Bundler/internal/formatting"
	"github.com/AkhilSharma90/AI-Code-Bundler/internal/kafka"
	"github.com/go-git/go-git/v5"
)

func extractProjectName(url string) string {
	re := regexp.MustCompile(`github\.com/[^/]+/([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "unknown-project"
	}
	return strings.ToLower(matches[1])
}

func main() {
	// Parse command line argument for repo URL
	repoURL := flag.String("repo", "", "GitHub repository URL to bundle")
	flag.Parse()

	if *repoURL == "" {
		log.Fatal("Please provide a GitHub repository URL using -repo flag")
	}

	// Extract project name for output file
	projectName := extractProjectName(*repoURL)
	outputDir := "output"
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.txt", projectName))

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create temp directory for cloning
	tempDir := filepath.Join(os.TempDir(), projectName)
	os.RemoveAll(tempDir) // Clean up any existing directory

	// Clone repository
	log.Printf("Cloning repository: %s", *repoURL)
	_, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      *repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	// Use the bundling logic from bundle.go
	var standardPrefixesToIgnore = []string{
		".", "codefuse", "go", "license", "readme", "README",
		"pyproject.toml", "poetry.lock", "venv", "build", "dist", "out", "target", "bin",
		"node_modules", "coverage", "public", "static", "Thumbs.db", "package", "yarn.lock",
		"tsconfig", "next.config", "next-env", "__pycache__", "logs", "gradle", "CMakeLists",
		"vendor", "Gemfile", "composer", "tailwind", "postcss",
	}

	var standardExtensionsToIgnore = []string{
		".jpeg", ".jpg", ".png", ".gif", ".pdf", ".svg", ".ico", ".woff", ".woff2",
		".eot", ".ttf", ".otf",
	}
	extensionsToInclude := []string{}

	filePaths, err := files.GetAllFilePaths(tempDir, standardPrefixesToIgnore, extensionsToInclude, standardExtensionsToIgnore)
	if err != nil {
		log.Fatalf("Error retrieving file paths: %v", err)
	}

	projectTree := formatting.GeneratePathTree(filePaths)
	fileContentMap, err := files.GetContentMapOfFiles(filePaths, 100)
	if err != nil {
		log.Fatalf("Error getting file contents: %v", err)
	}

	projectString := formatting.CreateProjectString(projectTree, fileContentMap)

	// Save the bundled code to the output file
	if err := os.WriteFile(outputFile, []byte(projectString), 0644); err != nil {
		log.Fatalf("Error saving bundle: %v", err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(outputFile)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	log.Printf("Bundle saved to: %s", absPath)

	// Initialize Kafka producer
	producer, err := kafka.NewKafkaProducer("localhost:9092", "bundler-client", "code-bundles")
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Send file path to Kafka
	err = producer.Send(absPath)
	if err != nil {
		log.Fatalf("Failed to send to Kafka: %v", err)
	}

	log.Printf("Successfully sent bundle path to Kafka")

	// Clean up
	os.RemoveAll(tempDir)
}
