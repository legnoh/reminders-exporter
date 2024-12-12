package reminder

import (
	"encoding/json"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

type ReminderData struct {
	Reminders []Reminder
	Lists     []string
}

type Condition struct {
	Deadline  time.Duration `default:"-"`
	Completed bool          `default:"false"`
	List      []string      `default:"[]"`
	Notes     []string      `default:"[]"`
	Priority  []string      `default:"[]"`
	Title     []string      `default:"[]"`
}

type Filter struct {
	Name      string `default:"-"`
	Condition Condition
}

const (
	ReminderCLIPath  string = "/opt/homebrew/bin/reminders"
	PriorityLow      int    = 1
	PriorityMedium   int    = 5
	PriorityHigh     int    = 9
	PriorityNoneReal int    = 0
)

var (
	Priorities = map[int]string{
		0: "none",
		1: "low",
		5: "medium",
		9: "high",
	}
	log *logrus.Logger
)

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func FetchReminderData() ReminderData {

	var reminders []Reminder
	jsonReminders, err := exec.Command(ReminderCLIPath, "show-all", "--include-completed", "--format=json").Output()
	if err != nil {
		log.Fatalf("Fetch Reminder Data failed. Please check your permission in 'Privacy and Security' setting!")
	}
	if json.Unmarshal(jsonReminders, &reminders) != nil {
		log.Fatalf("Unmarshal Reminder Data failed. reminders-cli is working collect?")
	}

	var lists []string
	jsonLists, err := exec.Command(ReminderCLIPath, "show-lists", "--format=json").Output()
	if err != nil {
		log.Fatalf("Fetch List Data failed. Please check your permission in 'Privacy and Security' setting!")
	}
	if json.Unmarshal(jsonLists, &lists) != nil {
		log.Fatalf("Unmarshal List Data failed. reminders-cli is working collect?")
	}

	return ReminderData{
		Reminders: reminders,
		Lists:     lists,
	}
}

func (reminders ReminderData) Count(condition Condition) float64 {
	matchedReminders := make([]Reminder, 0)
	for _, reminder := range reminders.Reminders {
		if condition.Deadline != 0 {
			condDueDate := time.Now().Add(condition.Deadline)
			if reminder.DueDate.IsZero() || condDueDate.Before(reminder.DueDate) {
				continue
			}
		}
		if condition.Completed && !reminder.IsCompleted {
			continue
		}
		if !condition.Completed && reminder.IsCompleted {
			continue
		}
		if len(condition.List) != 0 {
			if !slices.Contains(condition.List, reminder.List) {
				continue
			}
		}
		if len(condition.Priority) != 0 {
			if !slices.Contains(condition.Priority, Priorities[reminder.Priority]) {
				continue
			}
		}
		if len(condition.Title) != 0 {
			var hit bool = false
			for _, keyword := range condition.Title {
				if strings.Contains(reminder.Title, keyword) {
					hit = true
				}
			}
			if !hit {
				continue
			}
		}
		if len(condition.Notes) != 0 {
			var hit bool = false
			for _, keyword := range condition.Notes {
				if strings.Contains(reminder.Notes, keyword) {
					hit = true
				}
			}
			if !hit {
				continue
			}
		}
		matchedReminders = append(matchedReminders, reminder)
	}
	return float64(len(matchedReminders))
}
