package models

type User struct {
	Name       string  `toml:"name" json:"name" xml:"name"`
	Age        int     `toml:"age" json:"age" xml:"age"`
	Salary     float32 `toml:"salary" json:"salary" xml:"salary"`
	Occupation string  `toml:"occupation" json:"occupation" xml:"occupation"`
}
