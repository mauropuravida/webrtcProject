package db


type CameraDAO struct {
}

func InsertCam(w http.ResponseWriter, r *http.Request) error {
	query := "INSERT INTO Cameras(user, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?)"
	db := get()
	defer db.Close()
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	user := r.FormValue("user")
	


	defer stmt.Close()
	/*result, err := stmt.Exec(user, false, time.now() , u.loc, u.token_session_camera, u.token_session_consumer, u.id_camera)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
*/
	)
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

*/
func deleteCam(w http.ResponseWriter, r *http.Request) error {
	query := "DELETE FROM Cameras WHERE user=?"
	
	db := get()
	defer db.Close()

	user:= r.URL.Query().Get("id")

	log.Println("delete"+user)
	stmt, err := db.Prepare(query)
	if err != nil {
		return  err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(user)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", 301)
	
	return nil


}
