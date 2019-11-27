package global


type Mark struct {
	// db: Maps the "Identifier" property to the "identifier" column
	// of the "Mark" table.
	// json: Maps the "Identifier" property to the "ident" record field
	// of the json record
	Identifier string `json:"ident" db:"identifier"`
	Email      string `json:"email" db:"email"`
	Firstname  string `json:"firstname" db:"firstname"`
	Lastname   string `json:"lastname" db:"lastname"`
	// Metadata   *Metadata `json:"metadata:optional" db:"metadata:optional"`
}

type Metadata struct {
	State string `json:"state"`
}
