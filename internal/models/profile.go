package models

type Profile struct {
	ICCID      string `db:"iccid" json:"iccid"`
	Status     string `db:"status" json:"status"`
	MatchingID string `db:"matching_id" json:"matching_id"`
}

func (p Profile) CVSRespond() []string {
	var csvSlice = make([]string, 0, 3)

	csvSlice = append(csvSlice, p.ICCID, p.MatchingID)

	return csvSlice
}
