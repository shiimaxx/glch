# glch

`glch` generates changelog based on GitLab Merge Request.


## Usage

Run the `glch` after you moved to the GitLab project root directory.

```
$ glch [--latest|--only version|--next-version version]
```

You can use the following options.

- `--latest`
    - Display changelog for the most recent version only
- `--only version`
    - Display changelog for the specified versin only
- `--next-version version`
    - Set the next version (default is "Unreleased")


## GitLab Token

Please set your GitLab Token when using in a private project. You can get GitLab Token from [this page](https://gitlab.com/profile/personal_access_tokens).

```
$ export GITLAB_TOKEN=...
```


## GitLab API Endpoint

Default GitLab API Endpoint is `https://gitlab.com/api/v4/`. You can change it via GITLAB_API.

```
export GITLAB_API=https://gitlab.example.com/api/v4/
```


## Format of the changelog

`glch` will generate a changelog following format that based on 「[keep a changelog](https://keepachangelog.com/)」

```
## Version - YYYY-MM-DD

- <Merge Request title> <Reference> from @<Author>
- ...
```


## Install

- Download binary from [release page](https://github.com/shiimaxx/glch/releases)
- Copy binary to `$PATH` directory
