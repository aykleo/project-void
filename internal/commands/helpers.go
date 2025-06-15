package commands

import (
	"fmt"
	"strings"
	"time"
)

func GetAuthorNamesFromCommand(commandName string) []string {
	if !strings.HasPrefix(commandName, "git a ") {
		return nil
	}

	authorPart := strings.TrimPrefix(commandName, "git a ")
	authorPart = strings.TrimSpace(authorPart)

	if authorPart == "" {
		return nil
	}

	authors := strings.Split(authorPart, ",")
	var cleanAuthors []string
	for _, author := range authors {
		cleanAuthor := strings.TrimSpace(author)
		if cleanAuthor != "" {
			cleanAuthors = append(cleanAuthors, cleanAuthor)
		}
	}

	return cleanAuthors
}

func GetDateFromCommand(commandName string) (time.Time, error) {
	if !strings.HasPrefix(commandName, "void sd ") {
		return time.Time{}, fmt.Errorf("not a void sd command")
	}

	datePart := strings.TrimPrefix(commandName, "void sd ")
	datePart = strings.TrimSpace(datePart)

	if datePart == "" {
		return time.Time{}, fmt.Errorf("no date provided")
	}

	parsedDate, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	return parsedDate, nil
}
