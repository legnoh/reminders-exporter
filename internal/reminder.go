package internal

import (
	"encoding/json"
	"os/exec"
	"time"
)

type Reminder struct {
	DueDate     time.Time `json:"dueDate,omitempty"`
	ExternalId  string    `json:"externalId,omitempty"`
	IsCompleted bool      `json:"isCompleted,omitempty"`
	List        string    `json:"list,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Title       string    `json:"title,omitempty"`
}

func GetAllReminder() []Reminder {
	var reminders []Reminder
	jsonList, err := exec.Command("/opt/homebrew/bin/reminders", "show-all", "--format=json").Output()
	if err == nil {
		json.Unmarshal(jsonList, &reminders)
	}
	return reminders
}

func GetAllList() []string {
	var lists []string
	jsonList, err := exec.Command("/opt/homebrew/bin/reminders", "show-lists", "--format=json").Output()
	if err == nil {
		json.Unmarshal(jsonList, &lists)
	}
	return lists
}

func SearchReminder(reminders []Reminder, conditions Reminder) []Reminder {
	var matchedReminders []Reminder
	for _, reminder := range reminders {
		if !conditions.DueDate.IsZero() {
			if reminder.DueDate.IsZero() || conditions.DueDate.Before(reminder.DueDate) {
				continue
			}
		}
		matchedReminders = append(matchedReminders, reminder)
	}
	return matchedReminders
}
