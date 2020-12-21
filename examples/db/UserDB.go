package db
import ("database/sql"
	_"project/webrtcProject/examples/models"
)
/*
func GetAll() ([]models.user, error) {
	query := "SELECT * FROM Users"
	users := make([]models.user, 0)
	db := get()
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var row models.user
		err := rows.Scan(&row.ID, &row.FirstName, &row.LastName, &row.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, row)
	}

	return users, nil

}*/
func Delete(id int) ( sql.Result,error) {
	query := "DELETE FROM Users WHERE id=?"
	db := get()
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return  nil, err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(id)
	if err != nil {
		return nil, err
	}

	
	return rows, nil

}