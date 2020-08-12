package batch

// Request data model for an incoming request
type Request struct {
	UserID    int      `json:"user_id"`
	Total     float32  `json:"total"`
	Title     string   `json:"title"`
	Meta      Metadata `json:"meta"`
	Completed bool     `json:"completed"`
}

// Metadata data model for request metadata
type Metadata struct {
	Logins  []Login      `json:"logins"`
	Numbers PhoneNumbers `json:"phone_numbers"`
}

// Login data model for request metadata login info
type Login struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

// PhoneNumbers data model for request metadata phone number info
type PhoneNumbers struct {
	Home   string `json:"home"`
	Mobile string `json:"mobile"`
}

// Config is a data model for the repositories configuration
// includes batch size, batch interval, and post endpoint
type Config struct {
	Size     int
	Interval int
	Endpoint string
}
