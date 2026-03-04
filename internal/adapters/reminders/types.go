package reminders

// ReminderJSON represents a single reminder as returned by JXA scripts.
type ReminderJSON struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Body             string  `json:"body"`
	DueDate          *string `json:"dueDate"`
	Priority         int     `json:"priority"`
	Completed        bool    `json:"completed"`
	Flagged          bool    `json:"flagged"`
	CreationDate     string  `json:"creationDate"`
	ModificationDate string  `json:"modificationDate"`
}
