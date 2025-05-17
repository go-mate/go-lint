<p align="center">
  <img 
    alt="golangci-lint logo" 
    src="assets/golangci-lint-logo.jpeg" 
    style="max-height: 500px; width: auto; max-width: 100%;" 
  />
</p>
<h3 align="center">golangci-lint</h3>
<p align="center">execute <code>golangci-lint run</code> with golang os/exec</p>

---

# go-lint
Execute `golangci-lint run` with Golang's `os/exec`.

# install

```bash
go install github.com/go-mate/go-lint/cmd/go-lint@latest
```

# command

#### Single Project
```bash
cd project-path && go-lint
```
Analyze and report lint issues for a Go project.

#### Multiple Subprojects
```bash
cd awesome-path && go-lint
```
Analyze and report lint issues for all Go subprojects within the given path.

**Supported version:**
```bash
golangci-lint version
```

Output:
```text
golangci-lint has version 2.0.2 built with go1.24.1 from 2b224c2 on 2025-03-25T20:33:26Z
```

## License

MIT License. See [LICENSE](LICENSE).

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

If you find this package valuable, give me some stars on GitHub! Thank you!!!

**Thank you for your support!**

**Happy Coding with `go-lint`!** 🎉

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/go-lint.svg?variant=adaptive)](https://starchart.cc/go-mate/go-lint)
