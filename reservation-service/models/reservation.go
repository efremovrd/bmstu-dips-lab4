package models

type Reservation struct {
	Id, Hotel_id                                                         int
	Username, Status, Reservation_uid, Payment_uid, Start_date, End_data string
}
