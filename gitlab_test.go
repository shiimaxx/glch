package main

import (
	"testing"

	"github.com/xanzy/go-gitlab"
)

// testProjectID is project id of https://gitlab.com/shiimaxx/glch-demo
var testProjectID = 18674337

func testGitlabClient(t *testing.T) *gitlabClient {
	t.Helper()

	c, err := gitlab.NewClient("", gitlab.WithBaseURL(defaultGitLabAPIEndpoint))
	if err != nil {
		t.Fatal("create gitlab client failed: ", err)
	}
	return &gitlabClient{c}
}

func Test_getCommits(t *testing.T) {
	 testClient := testGitlabClient(t)

	 wants := []string{
		 "Merge branch 'feature-3' into 'master'",
		 "Add feature 3",
		 "v0.2.0",
		 "Merge branch 'feature-2' into 'master'",
		 "Add feature 2",
		 "Merge branch 'feature-1' into 'master'",
		 "Add feature 1",
		 "v0.1.0",
		 "Add",
		 "Initial commit",
	 }

	 commits, err := testClient.getCommits(testProjectID)
	 if err != nil {
		 t.Fatal("get commits failed: ", err)
	 }

	for i, commit := range commits {
		got := commit.Title
		want := wants[i]
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func Test_getMergeRequests(t *testing.T) {
	 testClient := testGitlabClient(t)

	 wants := []string{
		 "Feature 3",
		 "Feature 2",
		 "Feature 1",
	 }

	 mergeRequests, err := testClient.getMergeRequest(testProjectID)
	 if err != nil {
		 t.Fatal("get merge requests failed: ", err)
	 }

	for i, mr := range mergeRequests {
		got := mr.Title
		want := wants[i]
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func Test_getProject(t *testing.T) {
	 testClient := testGitlabClient(t)

	 if _, err := testClient.getProject("shiimaxx/glch-demo"); err != nil {
		 t.Fatal("get project failed: ", err)
	 }
}

func Test_getTags(t *testing.T) {
	 testClient := testGitlabClient(t)

	 wants := []string{
		 "v0.2.0",
		 "v0.1.0",
	 }

	 tags, err := testClient.getTags(testProjectID)
	 if err != nil {
		 t.Fatal("get tags failed: ", err)
	 }

	for i, tag := range tags {
		got := tag.Name	
		want := wants[i]
		if got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}

}


