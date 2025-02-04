package llm

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func fetchReviewString() string {
	file, err := os.Open("./data/sample_reviews.txt")
	if err != nil {
		log.Println("Error reading file contents=", err)
		return ""
	}

	var reviewString strings.Builder

	var reviews []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reviews = append(reviews, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	contructBulletPoints(reviews, &reviewString)
	return reviewString.String()
}

func contructBulletPoints(reviews []string, reviewString *strings.Builder) {
	for i, review := range reviews {
		reviewString.WriteString(strconv.Itoa(i+1) + ". " + review + "\n")
	}
}
