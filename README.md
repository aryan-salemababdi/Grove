# 🌱 Grove (v0.2.1)
A minimal **Modular framework skeleton** for Go.

Grove is an experimental framework inspired by [NestJS](https://nestjs.com/), built with:
- ⚡️ [Fiber](https://github.com/gofiber/fiber) for fast HTTP.
- 🧩 [Uber Dig](https://github.com/uber-go/dig) for dependency injection.
- 🛠️ [Cobra](https://github.com/spf13/cobra) for CLI commands.

---

## 🚀 Getting Started

### Install
Make sure you have Go 1.21+ installed.

```bash
go install github.com/aryan-salemababdi/Grove@latest
```

## ⚠️ Ensure your Go bin path is in $PATH:

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
(Add it permanently in ~/.zshrc or ~/.bashrc)
```

## Create a new app

```bash
grove new myapp
cd myapp
go run main.go
```

Now open http://localhost:3000 → you should see:

```bash
hello from Grove!
```


## 🧩 CLI Commands

	•	grove new <name> → creates a new Grove project with a default module.
	•	grove g module <name> → generates a scaffolded module inside your project.



## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first.
