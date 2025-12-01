package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type Habit struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"complete"`
}

var habits = make(map[int]Habit)
var nextID = 1

func addHabit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var h Habit
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.ID = nextID
	nextID++
	habits[h.ID] = h
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(h)
	// ignore error for now
	if err != nil {
		return
	}
	habits[h.ID] = h
	saveHabitsToFile()
}

func listHabits(w http.ResponseWriter, r *http.Request) {
	var habitList []Habit
	for _, h := range habits {
		habitList = append(habitList, h)
	}
	err := json.NewEncoder(w).Encode(habitList)
	if err != nil {
		return
	}
}

func toggleHabit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	habit, exists := habits[id]
	if !exists {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	habit.Completed = !habit.Completed
	habits[id] = habit

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habit)
	err = saveHabitsToFile()
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(habit)
}

func main() {

	if err := loadHabitsFromFile(); err != nil {
		// for now just log to console
		println("could not load habits:", err.Error())
	}
	// Serve the static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/toggle", toggleHabit)
	// Keep your /habits API handlers as they are
	http.HandleFunc("/habits", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listHabits(w, r)
		} else if r.Method == http.MethodPost {
			addHabit(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

const dataFile = "habits.json"

func saveHabitsToFile() error {
	f, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	enc := json.NewEncoder(f)
	return enc.Encode(habits)
}

func loadHabitsFromFile() error {
	f, err := os.Open(dataFile)
	if err != nil {
		// First run: file may not exist yet
		return nil
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&habits); err != nil {
		return err
	}

	// Recalculate nextID so new habits get unique IDs
	maxID := 0
	for id := range habits {
		if id > maxID {
			maxID = id
		}
	}
	nextID = maxID + 1
	return nil
}
