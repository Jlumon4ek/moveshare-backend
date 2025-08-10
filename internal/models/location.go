package models

type State struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type City struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	StateID int64  `json:"state_id" db:"state_id"`
}

type CityWithState struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StateID   int64  `json:"state_id"`
	StateName string `json:"state_name"`
}

type CitiesQuery struct {
	StateID *int64 `form:"state_id"`
}
