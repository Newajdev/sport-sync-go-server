package dto

type Response struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	ZoneID       uint   `json:"zone_id"`
	LicensePlate string `json:"license_plate"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type MyReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
}

type AdminResponse struct {
	ID           uint        `json:"id"`
	UserID       uint        `json:"user_id"`
	ZoneID       uint        `json:"zone_id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	User         UserSummary `json:"user"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
}
