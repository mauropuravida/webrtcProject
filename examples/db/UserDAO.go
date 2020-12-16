package db


type UserDAO struct {
}

func (dao UserDAO) GetAll() ([]models.User, error) {
	query := "SELECT * FROM Users"
	users := make([]models.User, 0)
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
		var row models.User
		err := rows.Scan(&row.ID, &row.FirstName, &row.LastName, &row.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, row)
	}

	return users, nil

}
func (dao UserDAO) Delete(id int) ([]models.User, error) {
	query := "DELETE FROM Users WHERE id=?"
	users := make([]models.User, 0)
	db := get()
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return  err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(id)
	if err != nil {
		return err
	}

	
	return nil

}