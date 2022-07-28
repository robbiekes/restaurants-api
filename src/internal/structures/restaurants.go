package structures

type Info struct {
	Name     string       `json:"name"`
	Total    int          `json:"total"`
	Rests    []Restaurant `json:"restaurants"`
	PrevPage int          `json:"prev_page"`
	NextPage int          `json:"next_page"`
	LastPage int64        `json:"last_page"`
}

type ClosestRests struct {
	Name  string       `json:"name"`
	Rests []Restaurant `json:"restaurants"`
}

type Restaurant struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
