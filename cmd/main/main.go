package main

import (
	"fmt"
	"os"

	// "time"

	// git "project-void/internal/git"
	// hello_world "project-void/internal/ui/hello-world"
	"project-void/internal/ui/home"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := home.InitialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// for {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: program <git-repo-path> <since-date>")
// 		fmt.Println("Please enter a valid path:")
// 		var path string
// 		fmt.Scanln(&path)
// 		os.Args = append(os.Args, path)
// 		continue
// 	}

// 	repoPath := os.Args[1]
// 	sinceDate, err := time.Parse("02/01/2006", os.Args[2])
// 	if err != nil {
// 		fmt.Printf("Error parsing date: %v\n", err)
// 		fmt.Println("Please enter a valid date (dd/mm/yyyy format):")
// 		var date string
// 		fmt.Scanln(&date)
// 		os.Args[2] = date
// 		continue
// 	}

// 	commits, err := git.GetCommitsSince(repoPath, sinceDate)
// 	if err != nil {
// 		fmt.Printf("Error getting commits: %v\n", err)
// 		fmt.Println("Please enter a valid git repository path:")
// 		var path string
// 		fmt.Scanln(&path)
// 		os.Args[1] = path
// 		continue
// 	}

// 	p := tea.NewProgram(hello_world.InitialModel(commits))
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Error: %v", err)
// 		os.Exit(1)
// 	}
// 	break
// }
