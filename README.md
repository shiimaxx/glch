# glch

`glch` generates changelog based on GitLab Merge Request.


## Usage

Run the `glch` after you moved to GitLab project root directory.

```
$ glch [--latest|--only version|--next-version version]
```

You can use following options.

- `--latest`
    - Display changelog for the most recent version only
- `--only version`
    - Display changelog for the specified versin only
- `--next-version version`
    - Set the next version (default is "Unreleased")


## GitLab Token

Setting your GitLab Token. You can get GitLab Token from [this page](https://gitlab.com/profile/personal_access_tokens).

```
$ export GITLAB_TOKEN=...
```


## GitLab API Endpoint

Default GitLab API Endpoint is `https://gitlab.com/api/v4/`. You can change it via GITLAB_API.

```
export GITLAB_API=https://gitlab.example.com/api/v4/
```


## Format of the changelog

`glch` will generate changelog following format.

```
## Version - YYYY-MM-DD

- <Merge Request title> <Reference> from @<Author>
- ...
```


## Install

- Download binary from [release page](https://github.com/shiimaxx/glch/releases)
- Copy binary to `$PATH` directory
