package main

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type Bag struct {
	ID          int
	UserID      int
	RoasterName string
	BeanOrigin  string
	RoastDate   string
}

type Shot struct {
	ID              int
	BagId           int
	GrindSize       float64
	DoseGrams       float64
	YeildGrams      float64
	BrewTimeSeconds int
	Taste           string
	Rating          int
	IsDialedIn      bool
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateBagRequest struct {
	RoasterName string `json:"roasterName"`
	BeanOrigin  string `json:"beanOrigin"`
	RoastDate   string `json:"roastDate"`
}
