package models 

import ("time"
)

type Camera struct{
	ID int
	User_id int
	Active int
	Url string
	Created time.Time
	Loc string
	T_s_cam string
	T_s_con string
	Id_cam int
	
}

