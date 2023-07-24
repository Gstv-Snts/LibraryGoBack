package models

type User struct {
	Id       int
	Email    string
	Password string
}

type Status struct {
	IsActive    bool
	Description string
}

type Rent struct {
	StudentName   string
	Class         string
	WithdrawlDate string
	DeliveryDate  string
}

type Book struct {
	Id              int
	Tittle          string
	Author          string
	Genre           string
	Status          Status
	Image           string
	SystemEntryDate string
	Synopsis        string
	RentHistory     []Rent
}

type BookUpdate struct {
	Tittle          string
	Author          string
	Genre           string
	Status          Status
	Image           string
	SystemEntryDate string
	Synopsis        string
}
