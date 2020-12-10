package mysql

import "models"

type CameraDAO struct {
}

func (dao CameraDAO) Create(u *models.Camera) error {
	query := "INSERT INTO Cameras(user, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?)"
	db := get()
	defer db.Close()
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(u.user, u.active, u.created, u.loc, u.token_session_camera, u.token_session_consumer, u.id_camera)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)
	return nil
}
/*
func (dao UserImplMysql) GetAll() ([]models.User, error) {
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
func (dao UserImplMysql) Delete(id int) ([]models.User, error) {
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
*/