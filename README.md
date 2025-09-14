# 🌱 Velora (v0.4.4)
A minimal **Modular framework skeleton** for Go.

Velora is an experimental framework inspired by [NestJS](https://nestjs.com/), built with:
- ⚡️ [Fiber](https://github.com/gofiber/fiber) for fast HTTP.
- 🧩 [Uber Dig](https://github.com/uber-go/dig) for dependency injection.
- 🛠️ [Cobra](https://github.com/spf13/cobra) for CLI commands.

---

## 🚀 Getting Started

### Install
Make sure you have Go 1.21+ installed.

```bash
go install github.com/aryan-salemababdi/Velora@latest
```

## ⚠️ Ensure your Go bin path is in $PATH:

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
(Add it permanently in ~/.zshrc or ~/.bashrc)
```

## Create a new app

```bash
velora new myapp
cd myapp
go run main.go
```

Now open http://localhost:3000 → you should see:

```bash
hello from Velora!
```


## 🧩 CLI Commands

	•	velora new <name> → creates a new Velora project with a default module.
	•	velora g module <name> → generates a scaffolded module inside your project.



## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first.
