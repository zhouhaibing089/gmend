package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// gmend provides an easy way:
//  1. go back to a history commit
//  2. make changes
//  3. restore to HEAD
// usage:
//  gmend <commit>

// usage prints information on how to use this
func usage() {
	fmt.Printf(`%s is a cli tool helps you to amend commits in an easy way.

usage:
	%s <commit>

commit: the commit id you would like to amend.
`, os.Args[0], os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	// make sure that we are working on a git repository
	if isGitRepo() != nil {
		fmt.Println("Not a Git Repository.")
		os.Exit(128)
	}
	// make sure the commit id present
	commit := os.Args[1]
	if isGitCommit(commit) != nil {
		fmt.Printf("<%s> Not Found\n", commit)
	}
	// make a backup for current branch
	branch, err := branchName()
	if err != nil {
		fmt.Println("failed to detect branch name")
		os.Exit(128)
	}
	branchName := fmt.Sprintf("%s-gmend-backup", branch)
	fmt.Printf("make a backup branch %s\n", branchName)
	if err := newBranch(branchName); err != nil {
		fmt.Printf("failed to create backup branch %q", branchName)
		os.Exit(128)
	}
	// get all the commit ids from the one that we want to amend
	commits, err := listCommits(commit)
	if err != nil {
		fmt.Printf("failed to list commits since %s\n", commit)
		os.Exit(128)
	}
	// reset to the commit that we would like to change
	fmt.Printf("Reset to %s\n", commit)
	if err := resetTo(commit); err != nil {
		fmt.Printf("failed to checkout %s\n", commit)
		os.Exit(128)
	}

	// now make your changes and wait for changes
	fmt.Printf("Now make your changes, once ready, press Enter")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
	fmt.Printf("Save your changes: ")
	if err := save(); err != nil {
		fmt.Println("failed")
		os.Exit(128)
	} else {
		fmt.Println("ok")
	}
	// reapply back
	for i := len(commits) - 1; i >= 0; i-- {
		target := commits[i]
		fmt.Printf("Apply %s:", target)
		if err := applyCommit(target); err != nil {
			fmt.Println("failed")
			os.Exit(128)
		} else {
			fmt.Println("ok")
		}
	}
}

func isGitRepo() error {
	cmd := exec.Command("git", "status")
	return cmd.Run()
}

func isGitCommit(commit string) error {
	cmd := exec.Command("git", "rev-parse", commit)
	return cmd.Run()
}

func branchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func newBranch(name string) error {
	cmd := exec.Command("git", "branch", name)
	return cmd.Run()
}

func deleteBranch(name string) error {
	cmd := exec.Command("git", "branch", "-D", name)
	return cmd.Run()
}

func listCommits(since string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", since+"..HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func resetTo(commit string) error {
	cmd := exec.Command("git", "reset", commit, "--hard")
	return cmd.Run()
}

func applyCommit(commit string) error {
	cmd := exec.Command("git", "cherry-pick", commit)
	return cmd.Run()
}

func save() error {
	cmd := exec.Command("git", "commit", "-a", "--amend", "--no-edit")
	return cmd.Run()
}
