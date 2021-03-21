package main

import (
	"bytes"
	"testing"
)

func Test_run(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)

	glch := glch{
		repository: testGitlabClient(t),
		projectPath: "shiimaxx/glch-demo",

		outStream: outStream,
		errStream: errStream,
	}
	args := []string{"glch"}

	exitCode := glch.run(args)
	if exitCode != 0 {
		t.Fatal(errStream.String())
	}

	want := `## Unreleased - 2021-03-21

- Feature 3 shiimaxx/glch-demo!3 from [@shiimaxx](https://gitlab.com/shiimaxx)


## v0.2.0 - 2020-05-09

- Feature 2 shiimaxx/glch-demo!2 from [@shiimaxx](https://gitlab.com/shiimaxx)
- Feature 1 shiimaxx/glch-demo!1 from [@shiimaxx](https://gitlab.com/shiimaxx)


## v0.1.0 - 2020-05-09

- Initial release


`

	got := outStream.String()

	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
