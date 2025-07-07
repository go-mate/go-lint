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
Analyze and report lint issues for the Go project.

#### Multiple Subprojects
```bash
cd awesome-path && go-lint
```
Analyze and report lint issues for Go subprojects.

```bash
cd awesome-path && go-lint run --debug=0
```

```bash
cd awesome-path && go-lint run --debug=1
```

**Supported version:**
```bash
golangci-lint version
```

Output:
```text
golangci-lint has version 2.0.2 built with go1.24.1 from 2b224c2 on 2025-03-25T20:33:26Z
```

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repo on GitHub (using the webpage interface).
2. Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. Navigate to the cloned project (`cd repo-name`)
4. Create a feature branch (`git checkout -b feature/xxx`).
5. Stage changes (`git add .`)
6. Commit changes (`git commit -m "Add feature xxx"`).
7. Push to the branch (`git push origin feature/xxx`).
8. Open a pull request on GitHub (on the GitHub webpage).

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

If you find this package valuable, give me some stars on GitHub! Thank you!!!

**Thank you for your support!**

**Happy Coding with this package!** 🎉

Give me stars. Thank you!!!

---

## GitHub Stars

[![starring](https://starchart.cc/go-mate/go-lint.svg?variant=adaptive)](https://starchart.cc/go-mate/go-lint)
