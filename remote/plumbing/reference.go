package plumbing

import (
	"fmt"
	"regexp"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	IssueBranchPrefix        = "issues"
	MergeRequestBranchPrefix = "merges"
)

// isBranch checks whether a reference name indicates a branch
func IsBranch(name string) bool {
	return plumbing.ReferenceName(name).IsBranch()
}

// IsIssueReference checks whether a branch is an issue branch
func IsIssueReference(name string) bool {
	return regexp.MustCompile(fmt.Sprintf("^refs/heads/%s/[1-9]+([0-9]+)?$", IssueBranchPrefix)).MatchString(name)
}

// IsIssueReferencePath checks if the specified reference matches an issue reference path
func IsIssueReferencePath(name string) bool {
	return regexp.MustCompile(fmt.Sprintf("^refs/heads/%s(/|$)?", IssueBranchPrefix)).MatchString(name)
}

// IsMergeRequestReference checks whether a branch is a merge request branch
func IsMergeRequestReference(name string) bool {
	re := "^refs/heads/%s/[1-9]+([0-9]+)?$"
	return regexp.MustCompile(fmt.Sprintf(re, MergeRequestBranchPrefix)).MatchString(name)
}

// IsMergeRequestReferencePath checks if the specified reference matches a merge request reference path
func IsMergeRequestReferencePath(name string) bool {
	re := "^refs/heads/%s(/|$)?"
	return regexp.MustCompile(fmt.Sprintf(re, MergeRequestBranchPrefix)).MatchString(name)
}

// isReference checks the given name is a reference path or full reference name
func IsReference(name string) bool {
	re := "^refs/(heads|tags|notes)((/[a-z0-9_-]+)+)?$"
	return regexp.MustCompile(re).MatchString(name)
}

// isTag checks whether a reference name indicates a tag
func IsTag(name string) bool {
	return plumbing.ReferenceName(name).IsTag()
}

// isNote checks whether a reference name indicates a tag
func IsNote(name string) bool {
	return plumbing.ReferenceName(name).IsNote()
}

// MakeIssueReference creates an issue reference
func MakeIssueReference(id interface{}) string {
	return fmt.Sprintf("refs/heads/%s/%v", IssueBranchPrefix, id)
}

// MakeIssueReferencePath returns the full issue reference path
func MakeIssueReferencePath() string {
	return fmt.Sprintf("refs/heads/%s", IssueBranchPrefix)
}

// MakeMergeRequestReference creates a merge request reference
func MakeMergeRequestReference(id interface{}) string {
	return fmt.Sprintf("refs/heads/%s/%v", MergeRequestBranchPrefix, id)
}

// MakeIssueReferencePath returns the full merge request reference path
func MakeMergeRequestReferencePath() string {
	return fmt.Sprintf("refs/heads/%s", MergeRequestBranchPrefix)
}
