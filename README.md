# poco-github-action
Github action for [poco](https://github.com/mudler/poco)


## Usage

It takes the same arguments as of poco, given you have a Dockerfile in the root of your repo:

```yaml
- name: Check out repository
  uses: actions/checkout@v2
- name: Build
  run: |
        docker build -t example ./
- name: Build a static binary out of it!
  uses: mudler/poco-github-action@main
  with:
    appDescription: "foo"
    appName: "bar"
    image: "example"
    appEntrypoint: "/usr/bin/wget"
    output: "./wget"
    appMounts: ""
    appStore: "$HOME/.store"
```
