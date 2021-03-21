package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

type doer interface {
	getCommits(int) ([]*gitlab.Commit, error)
	getMergeRequest(int) ([]*gitlab.MergeRequest, error)
	getProject(string) (*gitlab.Project, error)
	getTags(int) ([]*gitlab.Tag, error)
}

type gitlabClient struct {
	*gitlab.Client
}

func (c *gitlabClient) getTotalPages(endpoint string) (int, error) {
	req, err := c.NewRequest(
		"HEAD",
		fmt.Sprintf(endpoint),
		nil,
		nil,
	)
	if err != nil {
		return 0, errors.Wrap(err, "create new request failed")
	}

	res, err := c.Do(req, nil)
	if err != nil {
		return 0, errors.Wrap(err, "do request failed")
	}

	if res.StatusCode != http.StatusOK {
		return 0, errors.Errorf("api request failed: invalid status: %d", res.StatusCode)
	}

	xTotalPages := res.Header.Get("X-Total-Pages")
	totalPages, err := strconv.Atoi(xTotalPages)
	if err != nil {
		return 0, errors.Wrap(err, "convert X-Total-Pages failed")
	}

	return totalPages, nil
}

func (c *gitlabClient) getCommits(pid int) ([]*gitlab.Commit, error) {
	var commits []*gitlab.Commit
	for i := 1; ; i ++ {
		c, res, err := c.Commits.ListCommits(pid, &gitlab.ListCommitsOptions{
			ListOptions: gitlab.ListOptions{
				Page: i,
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "api requests failed")
		}
		if res.StatusCode != http.StatusOK {
			return nil, errors.Errorf("api requests failed: invalid status: %d", res.StatusCode)
		}
		commits = append(commits, c...)

		if res.NextPage == 0 {
			break
		}
	}

	return commits, nil
}

func (c *gitlabClient) getMergeRequest(pid int) ([]*gitlab.MergeRequest, error) {
	totalPages, err := c.getTotalPages(fmt.Sprintf("projects/%d/merge_requests", pid))
	if err != nil {
		return nil, errors.Wrap(err, "get total pages failed")
	}

	eg := errgroup.Group{}

	results := make([][]*gitlab.MergeRequest, totalPages)
	for i := 1; i <= totalPages; i++ {
		i := i
		eg.Go(func() error {
			m, res, err := c.MergeRequests.ListProjectMergeRequests(pid, &gitlab.ListProjectMergeRequestsOptions{
				ListOptions: gitlab.ListOptions{
					Page: i,
				},
			})
			if err != nil {
				return errors.Wrap(err, "api request failed")
			}
			if res.StatusCode != http.StatusOK {
				return errors.Errorf("api request failed: invalid status: %d", res.StatusCode)
			}
			results[i-1] = m
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var mergeRequests []*gitlab.MergeRequest
	for _, r := range results {
		mergeRequests = append(mergeRequests, r...)
	}

	return mergeRequests, nil
}

func (c *gitlabClient) getProject(path string) (*gitlab.Project, error) {
	project, res, err := c.Projects.GetProject(path, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get project failed")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("api request failed: invalid status: %d", res.StatusCode)
	}

	return project, nil
}

func (c *gitlabClient) getTags(pid int) ([]*gitlab.Tag, error) {
	totalPages, err := c.getTotalPages(fmt.Sprintf("projects/%d/repository/tags", pid))
	if err != nil {
		return nil, errors.Wrap(err, "get total pages failed")
	}

	eg := errgroup.Group{}

	results := make([][]*gitlab.Tag, totalPages)
	for i := 1; i <= totalPages; i++ {
		i := i
		eg.Go(func() error {
			t, res, err := c.Tags.ListTags(pid, &gitlab.ListTagsOptions{
				ListOptions: gitlab.ListOptions{
					Page: i,
				},
			})
			if err != nil {
				return errors.Wrap(err, "api request failed")
			}
			if res.StatusCode != http.StatusOK {
				return errors.Errorf("api request failed: invalid status: %d", res.StatusCode)
			}
			results[i-1] = t
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var tags []*gitlab.Tag
	for _, r := range results {
		tags = append(tags, r...)
	}

	return tags, nil
}
