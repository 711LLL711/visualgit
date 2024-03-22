package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func Scan(path string) {
	//scan for .git
	repo, err := scanGitFolder(path)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(repo)
	saveFilePath := GetSaveFilePath()
	//save to file
	WriteRepoToFile(saveFilePath, repo)
}

func GetSaveFilePath() string {
	//get database filepath
	usr, err := user.Current()
	if err != nil {
		fmt.Print("Get user home dir failed\n")
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.visualgit"
	return dotFile
}
func WriteRepoToFile(path string, repos []string) {
	//no duplicate
	oldrepos := fileToSlice(path)
	newrepos := addNewRepo(oldrepos, repos)
	slicetoFile(newrepos, path)
}

// open databasefile
func openFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("File not exist\n")
			_, err = os.Create(path)
			if err != nil {
				fmt.Print("Create file failed\n")
				log.Fatal(err)
			}
		} else {
			fmt.Print("Open file failed\n")
			log.Fatal(err)
		}
	}
	return file
}

func fileToSlice(path string) []string {
	f := openFile(path)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var repos []string
	for scanner.Scan() {
		line := scanner.Text()
		repos = append(repos, line)
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
	return repos
}

// add new repo
func addNewRepo(repos []string, newRepos []string) []string {
	for _, repo := range newRepos {
		if !contains(repos, repo) {
			repos = append(repos, repo)
		}
	}
	return repos
}
func contains(repos []string, repo string) bool {
	for _, r := range repos {
		if r == repo {
			return true
		}
	}
	return false
}

func slicetoFile(repos []string, path string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(path, []byte(content), 0755)
}

func scanGitFolder(folder string) ([]string, error) {
	var repos []string
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
		}
		return nil
	})
	return repos, err
}
