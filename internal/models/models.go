package models

import "time"

type University struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	HasDormitory bool   `json:"has_dormitory"`
	HasMilitary  bool   `json:"has_military"`
	Rating       int    `json:"rating"`
	ImageKey     string `json:"image_key"`
	ImageURL     string `json:"image_url"`
}

type Faculty struct {
	ID           int    `json:"id"`
	UniversityID int    `json:"university_id"`
	Name         string `json:"name"`
}

type Direction struct {
	ID             int    `json:"id"`
	FacultyID      int    `json:"faculty_id"`
	Name           string `json:"name"`
	TuitionFee     int    `json:"tuition_fee"`
	MinScoreBudget int    `json:"min_score_budget"`
	MinScorePaid   int    `json:"min_score_paid"`
	BudgetPlaces   int    `json:"budget_places"`
	PaidPlaces     int    `json:"paid_places"`
}

type Subject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Favorite struct {
	ID           int `json:"id"`
	UserID       int `json:"user_id"`
	UniversityID int `json:"university_id"`
}

type Filters struct {
	City       string `json:"city"`
	Form       string `json:"form"`
	Budget     bool   `json:"budget"`
	Paid       bool   `json:"paid"`
	Military   bool   `json:"military"`
	Dormitory  bool   `json:"dormitory"`
	TotalScore int    `json:"totalScore"`
	Subjects   []int  `json:"subjects"`
}

type SearchHistory struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Filters     string    `json:"filters"`
	SearchedAt  time.Time `json:"searched_at"`
	FiltersData Filters   `json:"filters_data"`
}

type PageData struct {
	Universities     []University `json:"universities"`
	Faculties        []Faculty    `json:"faculties"`
	Directions       []Direction  `json:"directions"`
	Subjects         []Subject    `json:"subjects"`
	Errors           Errors       `json:"errors"`
	City             string       `json:"city"`
	Form             string       `json:"form"`
	TotalScore       int          `json:"total_score"`
	Dormitory        bool         `json:"dormitory"`
	Military         bool         `json:"military"`
	Budget           bool         `json:"budget"`
	Paid             bool         `json:"paid"`
	SelectedSubjects []int        `json:"selected_subjects"`
	Offset           int
	HasMore          bool
}

type Errors struct {
	CityError  string `json:"city_error"`
	ScoreError string `json:"score_error"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	AuthProvider string `json:"auth_provider"`
	ProviderID   string `json:"provider_id"`
}

type ComparisonData struct {
	ID                int
	Option            string
	Name              string
	City              string
	HasDormitory      bool
	HasMilitary       bool
	TotalBudgetPlaces int
	TotalPaidPlaces   int
	MinTuitionFee     int
	MinScoreBudget    int
	MinScorePaid      int
	ProgramsCount     int
	SpecialtiesCount  int
}

type Favorites struct {
	Universities []University `json:"universities"`
	Faculties    []Faculty    `json:"faculties"`
	Directions   []Direction  `json:"directions"`
	Option       string
}
