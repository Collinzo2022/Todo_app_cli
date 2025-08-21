package main // program entry package

import (
	"fmt"  // used for formatted I/O, such as printing to the console
	"encoding/json" // for converting Go structs to JSON and back
	"strings" // For working with text (joining, lowercasing, etc.)
	"os" // for interacting with the operating system (file I/O, args)
	"time" // for time handling and timestamps

)
// represents a Todo item stored in memory and in JSON
type Todo struct {
	ID  int              `json:"id"` // unique task ID
	Title     string        `json:"title"`  // description/title of the task
	Completed bool          `json:"completed"` // whether the task is completed
	CreatedAt time.Time     `json:"created_at"` // creation timestamp
	CompletedAt *time.Time  `json:"completed_at,omitempty"`  // Time the task was completed  so it's a pointer

}

// TodoList is a slice of Todo
type TodoList struct {
	NextID int		`json:"next_id"`// next unique ID assign; persisted so IDs stay unique
	Items  []Todo	`json:"items"`  // the slice containing todos
}

// ===== Helper error functions =====
//if there’s an error, print it and stop the program immediately.
func errorChecker(err error, msg string) {
	if err != nil {
		fmt.Println("❌", msg, err)
		os.Exit(1)
	}
}

// print a warning but continue running.
func errorWarner(err error, msg string) {
	if err != nil {
		fmt.Println("⚠️ ", msg, err)
	}
}
//converts a command-line string into an integer ID and stops if invalid.
func parseID(arg string) int {
    var id int
    _, err := fmt.Sscanf(arg, "%d", &id)
    errorChecker(err, "Invalid ID:")
    return id
}
 //checks if the right number of arguments were given; if not, prints usage and exits.
func requiredArgs(numarg int, usage string) {
    if len(os.Args) < numarg {
        fmt.Println("Usage:", usage)
        os.Exit(1)
    }
}

// Add creates a new todo and appends it to the TodoList
func (t *TodoList) add(title string) {

	// Build the new Todo value using the NextID counter
	todo := Todo {
		ID: t.NextID,  // Set ID.
		Title: title, // Set task name.
		Completed: false,  // New task stacks incomplete.
		CreatedAt: time.Now(), // Timestamp of creation

	}
	t.Items = append(t.Items, todo) // Append the new todo to the slice.
	// Increment NextID so we never reuse a numeric ID
	t.NextID++
	fmt.Println("✅ Todo added:", title)// confirm addition
}

// a function List that prints a friendly list view to stdout. Shows position, title, real ID, status and timestamps.
func (t *TodoList) List() {
	// If there are no items, print a message and return early.
	if len(t.Items) == 0 {
		fmt.Println("No todos yet!")
		return
	}
	//using for loop to loop through each Todo
	for i, todo := range t.Items {
		//default status for incomplete tasks
		status := "⏳ Pending"
		// if the task is done
		if todo.Completed {
			//update it to done
			status = "✅ Done"
		}
		// Print task details
		fmt.Printf("%d. %s | (ID:%d) | %s | Created: %s", i+1, todo.Title, todo.ID, status, todo.CreatedAt.Format("2006-01-02 15:04:05")) // Format time as YYYY-MM-DD HH:MM:SS.
		// If task is completed and has a completion time, print it.
	if todo.CompletedAt != nil {
		fmt.Printf(" | Finished: %s", todo.CompletedAt.Format("2006-01-02 15:04:05"))
	}
	//add new line for each task
	fmt.Println()


	}
}
// Complete method marks a task as done by searching for the real ID.
func (t *TodoList) Complete(id int) {
	// Iterate over the items by index so we can modify
	for i := range t.Items {
		// Find task with matching ID.
		if t.Items[i].ID == id {
			// If already completed...
			if t.Items[i].Completed {
				// Warn user.
				fmt.Println("Todo already completed.")
				//exit  function
				return
			}
			// Mark as done
			t.Items[i].Completed = true
			 // Get current time.
			now := time.Now()
			// Save completion time
			t.Items[i].CompletedAt = &now
			// print out the task
			fmt.Println(" 🎯Todo marked as completed:", t.Items[i].Title)
			//exit function
			return
		}
	}
	// todo not found If no matching ID found.
	fmt.Println("⚠️Todo not found. ")
}
// Edit method updates the title of a task.
func (t *TodoList) Edit(id int, newTitle string) {
	// Loop through tasks.
	for i := range t.Items {
		// Find task with matching ID.
		if t.Items[i].ID == id {
			// Change the title.
			t.Items[i].Title = newTitle
			// Confirmation.
			fmt.Println("✏️Todo updated.")
			// exit function
			return
		}
	}
	// If no matching ID not found
	fmt.Println("⚠️Todo not found.")
}
// Delete method removes a task from the list.
func (t *TodoList) Delete(id int) {
	// Loop through tasks.
	for i := range t.Items {
		// Find task with matching ID.
		if t.Items[i].ID == id {
			// Remove task from slice.
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			// Confirm deleted
			fmt.Println(" 🗑️ Todo deleted.")
			// exit function
			return
		}
	}
	// If no matching ID.
	fmt.Println("Todo not found.")
}
// Save writes the entire TodoList (NextID + Items) to the given filename as JSON.
func (t *TodoList) Save(filename string) error {
	file, err := os.Create(filename)
	// If conversion fails..
	errorChecker(err, "failed to save file")

	// Ensure file is closed when Save returns (success or error)
	defer file.Close()
	// Create a JSON encoder that writes directly to the file
	encoder := json.NewEncoder(file)
	// Make the output human-readable
	encoder.SetIndent("", " ")
	// Encode the structure to the file (streams bytes directly)
	errorChecker(encoder.Encode(t), "Failed to write to JSON:")
	return nil

}
// Load reads the TodoList from filename and populates t.
func (t *TodoList) Load(filename string) error {
	// Open file for reading
	file, err := os.Open(filename)
	// If error occurs
		// If file is missing, initialize an empty list and set NextID to 1
		if os.IsNotExist(err) {
			t.NextID = 1
			// Start with an empty list.
			t.Items = []Todo{}

			return nil
		}
		// Other errors are returned upward
	errorChecker(err, "failed to open file: ")
	// Load JSON into slice
	defer file.Close()
	// Decode JSON directly from the file into the struct
	decoder := json.NewDecoder(file)
	// If decode fails (malformed JSON), return the error
	errorChecker(decoder.Decode(t), "failed to decode JSON")


	// After loading, ensure NextID is at least max(existing IDs)+1 so we don't reuse IDs
		maxID := 0
		for _, todo := range t.Items {
			if todo.ID > maxID {
				maxID = todo.ID
			}
		}
	// if NextID was missing or <= maxID, advance it to maxID+1
		if t.NextID <= maxID {
			t.NextID = maxID + 1
		}
			// If NextID is still zero (e.g., file had no NextID field), ensure it is 1
		if t.NextID == 0 {
			t.NextID = 1
		}
		return nil

}


func main() {
	// File to save/load todos.
	const filename = "todos.json"
	// Create a TodoList variable.
	var todos TodoList

	// Load saved todos from file (if any).
	errorChecker(todos.Load(filename), "Error loading todos:")

	// Ensure user gave at least 1 argument (the command).
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [add|list|complete|edit|delete] ...")
		return
	}

	// Get the command and make it lowercase.
	command := strings.ToLower(os.Args[1])

	// Handle command.
	switch command {
	case "add":
		requiredArgs(3, "go run main.go add <title>")
		// Join all remaining arguments as title.
		title := strings.Join(os.Args[2:], " ")
		todos.add(title)
	case "list":
		// print out the list of todos
		todos.List()
	case "complete":
		requiredArgs(3, "go run main.go complete <id>")
		todos.Complete(parseID(os.Args[2]))

	case "edit":
		requiredArgs(4, "go run main.go edit <id> <new title>")
		todos.Edit(parseID(os.Args[2]), strings.Join(os.Args[3:], " "))
	case "delete":
		requiredArgs(3, "go run main.go delete <id>")
		todos.Delete(parseID(os.Args[2]))
	default:
		// Command not recognized.
		fmt.Println("⚠️Unknown command:", command)
	}
	// Save updated todos to file.
	errorWarner(todos.Save(filename), "Error saving todos")
}