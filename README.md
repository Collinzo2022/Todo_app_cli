# Todo App CLI

A simple command-line Todo application written in Go.
This app allows you to add, list, complete, edit, and delete your tasks directly from the terminal.
All tasks are saved in a local `todos.json` file for persistence.

---

## Features

- Add new todo items
- List all todos with status and timestamps
- Mark todos as completed
- Edit todo titles
- Delete todos
- Persistent storage using a JSON file

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or newer

### Installation

1. **Clone or download this repository.**
2. Open a terminal and navigate to the project directory:

   ```sh
   cd path/to/Todo_app_cli
   ```

3. **Run the application:**

   ```sh
   go run main.go [command] [arguments...]
   ```

---

## Usage

### Add a new todo

```sh
go run main.go add "Buy groceries"
```

### List all todos

```sh
go run main.go list
```

### Mark a todo as completed

```sh
go run main.go complete <id>
```

### Edit a todo

```sh
go run main.go edit <id> "New title"
```

### Delete a todo

```sh
go run main.go delete <id>
```

---

## Example

```sh
$ go run main.go add "Read Go documentation"
✅ Todo added: Read Go documentation

$ go run main.go list
1. Read Go documentation | (ID:1) | ⏳ Pending | Created: 2024-08-21 14:00:00

$ go run main.go complete 1
 🎯Todo marked as completed: Read Go documentation

$ go run main.go list
1. Read Go documentation | (ID:1) | ✅ Done | Created: 2024-08-21 14:00:00 | Finished: 2024-08-21 14:05:00
```

---

## Data Storage

- All todos are stored in a local `todos.json` file in the project