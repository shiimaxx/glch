package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tcnksm/go-gitconfig"
	"github.com/xanzy/go-gitlab"
)

const defaultGitLabAPIEndpoint = "https://gitlab.com/api/v4"

const (
	exitCodeOK int = 0

	exitCodeError = 10 + iota
	exitCodeParseError
	exitCodeInvalidResponseCode
)

var versionTagPattern = regexp.MustCompile(`v?[0-9]+\.[0-9]+\.[0-9]+$`)
var repoURLPattern = regexp.MustCompile(`([^/:]+)/([^/]+?)(?:\.git)?$`)

type glch struct {
	repository  doer
	projectPath string
	latest      bool
	only        string
	nextVersion string

	outStream, errStream io.Writer
}

type content struct {
	header string
	lines  []string
}

func (g *glch) run(args []string) int {
	var printVersion bool

	flags := flag.NewFlagSet("gch", flag.ContinueOnError)
	flags.SetOutput(g.outStream)
	flags.StringVar(&g.nextVersion, "next-version", "Unreleased", "")
	flags.StringVar(&g.only, "only", "", "")
	flags.BoolVar(&g.latest, "latest", false, "ignore when using -only option")
	flags.BoolVar(&printVersion, "version", false, "")
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(g.errStream, "flag parse failed: ", err)
		return exitCodeParseError
	}

	if printVersion {
		fmt.Fprintf(g.outStream, "%s version is %s\n", Name, Version)
		return exitCodeOK
	}

	project, err := g.repository.getProject(g.projectPath)
	if err != nil {
		fmt.Fprintln(g.errStream, "get project failed: ", err)
		return exitCodeError
	}

	tags, err := g.repository.getTags(project.ID)
	if err != nil {
		fmt.Fprintln(g.errStream, "get tags failed: ", err)
		return exitCodeError
	}

	var versionTags []*gitlab.Tag
	for _, tag := range tags {
		if versionTagPattern.MatchString(tag.Name) {
			versionTags = append(versionTags, tag)
		}
	}

	now := time.Now()
	versionTags = append(versionTags, &gitlab.Tag{
		Name: g.nextVersion,
		Commit: &gitlab.Commit{
			CreatedAt: &now,
		},
	})

	sort.Slice(versionTags, func(i, j int) bool {
		return versionTags[i].Commit.CreatedAt.After(*versionTags[j].Commit.CreatedAt)
	})

	commits, err := g.repository.getCommits(project.ID)
	if err != nil {
		fmt.Fprintln(g.errStream, "get commits failed: ", err)
		return exitCodeError
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].CreatedAt.After(*commits[j].CreatedAt)
	})

	version := g.nextVersion
	commitlog := make(map[string][]*gitlab.Commit)
	for _, commit := range commits {
		for _, t := range versionTags {
			if commit.ID == t.Commit.ID {
				version = t.Name
			}
		}
		if _, ok := commitlog[version]; !ok {
			commitlog[version] = []*gitlab.Commit{commit}
			continue
		}
		commitlog[version] = append(commitlog[version], commit)
	}

	mergeRequests, err := g.repository.getMergeRequest(project.ID)
	if err != nil {
		fmt.Fprintln(g.errStream, "get merge request failed: ", err)
		return exitCodeError
	}

	var changelog = make(map[string]*content)
	initialVersion := versionTags[len(versionTags)-1]
	latestVersion := versionTags[0]
	for i := 0; i < len(versionTags); i++ {
		var header string
		var li []string
		v := versionTags[i]

		header = fmt.Sprintf("## %s - %s", v.Name, v.Commit.CreatedAt.Format("2006-01-02"))
		for _, c := range commitlog[v.Name] {
			if strings.Contains(c.Message, "Merge branch") {
				for _, mr := range mergeRequests {
					if mr.MergeCommitSHA == c.ID {
						li = append(li, fmt.Sprintf("- %s %s!%d from [@%s](%s)",
							mr.Title,
							project.PathWithNamespace,
							mr.IID,
							mr.Author.Username,
							mr.Author.WebURL,
						))
					}
				}
			}
		}

		if v == initialVersion && len(li) < 1 {
			li = append(li, "- Initial release")
		}

		changelog[v.Name] = &content{
			header: header,
			lines:  li,
		}
	}

	ch := changelog[latestVersion.Name]
	if len(ch.lines) < 1 && len(versionTags) > 1 {
		versionTags = versionTags[1:]
	}

	if g.only != "" {
		for _, t := range versionTags {
			if g.only == t.Name {
				versionTags = []*gitlab.Tag{t}
				fmt.Fprint(g.outStream, display(versionTags, changelog))
				return exitCodeOK
			}
		}
		fmt.Fprintf(g.errStream, "%s is not found in version tags\n", g.only)
		return exitCodeError
	}

	if g.latest {
		versionTags = []*gitlab.Tag{versionTags[0]}
	}

	fmt.Fprint(g.outStream, display(versionTags, changelog))

	return exitCodeOK
}

func display(versionTags []*gitlab.Tag, changelog map[string]*content) string {
	var output string

	for _, v := range versionTags {
		c := changelog[v.Name]
		output += c.header + "\n\n"
		if len(c.lines) > 0 {
			output += strings.Join(c.lines, "\n")
			output += "\n\n"
		}
		output += "\n"
	}

	return output
}

func main() {
	endpoint := os.Getenv("GITLAB_API")
	if endpoint == "" {
		endpoint = defaultGitLabAPIEndpoint
	}

	gl, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(endpoint))
	if err != nil {
		log.Fatal("create gitlab client failed: ", err)
	}

	repo, err := gitconfig.OriginURL()
	if err != nil {
		log.Fatal("fetch origin url failed: ", err)
	}
	matches := repoURLPattern.FindStringSubmatch(repo)

	g := glch{
		repository:  &gitlabClient{gl},
		projectPath: fmt.Sprintf("%s/%s", matches[1], matches[2]),

		outStream: os.Stdout,
		errStream: os.Stderr,
	}

	os.Exit(g.run(os.Args))
}
