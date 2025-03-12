# gopm

Fast, efficient package manager for JavaScript and TypeScript

## Installation

For macOS, Linux, you can install gopm using the following command:

```bash
curl -sSL https://raw.githubusercontent.com/gopm/gopm/master/install.sh | sh
```

For Windows, you can download the installer from the [releases page](https://github.com/emmadal/gopm/releases).

## Usage

To use gopm, you can run the following command:

```bash
gopm <command>
```

For example, you can run:

```bash
gopm add <package> - Install a package and any packages that it depends on.
gopm install <package> - Install all packages from package.json.
gopm dev <package> - Install a package in development mode.
gopm rm <package> - Uninstall a package from node_modules and package.json.
gopm up <package> - Update a package in node_modules.
gopm list - List installed packages
gopm init - Initialize a new project
```

To show the help message, you can run:

```bash
gopm help
```

## Contributing

To contribute to gopm, please follow these steps:

1. Fork the repository on GitHub.
2. Clone your forked repository to your local machine.
3. Create a new branch for your changes.
4. Make your changes and commit them.
5. Push your changes to your forked repository.
6. Create a pull request on GitHub.

## License

gopm is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for more information.
