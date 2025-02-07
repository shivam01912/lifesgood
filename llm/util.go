package llm

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// TODO: remove this method once teh testing needs form a file are taken care of
func fetchReviewString() string {
	file, err := os.Open("./data/sample_reviews.txt")
	if err != nil {
		log.Println("Error reading file contents=", err)
		return ""
	}

	var reviews []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reviews = append(reviews, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return constructBulletPoints(reviews)
}

func constructBulletPoints(reviews []string) string {
	var reviewString strings.Builder

	for i, review := range reviews {
		reviewString.WriteString(strconv.Itoa(i+1) + ". " + review + "\n")
	}

	return reviewString.String()
}
