# CodeFuse CLI

`CodeFuse` is a command-line tool that allows you to easily bundle your codebase into a single file and get AI-powered code reviews. With CodeFuse, you can enhance code quality and catch bugs early, directly from your CLI.

## Features

- Bundle your codebase into a single `.txt` file.
- Get AI-driven code reviews on your bundled code.
- Customizable options (ignore/include specific files, directories, etc.).
- Written in Go for speed and efficiency.

## Important Commands

* **Bundle your codebase (saves locally as a `.txt` file)**:

```bash
   codefuse bundle
```

* **Get an AI code review of this bundle (saves locally as a .md file):

```bash
codefuse review
```