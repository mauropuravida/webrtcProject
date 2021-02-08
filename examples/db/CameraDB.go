package db

import (
	"time"

	"database/sql"
	"fmt"
	m "models"
)

func InsertCam(user int, loc string, url string, idcam int) error {

	query := "INSERT INTO Cameras(users_id, active, created, loc, url, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?,?)"

	db := get()

	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(user, false, time.Now(), loc, url, " ", " ", idcam)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

func UpdateCam(idCam int, user int, loc string, url string, active bool, tokencam string, tokencon string) (int64, error) {
	query := "UPDATE Cameras SET users_id=?, loc=?, url=?, active=?, token_session_camera=?, token_session_consumer=? WHERE id_camera=? and users_id=?"
	db := get()

	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(user, loc, url, active, tokencam, tokencon, idCam, user)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	defer stmt.Close()
	db.Close()
	return id, nil
}

func UpdateActiveCam(act bool, id int) (int64, error) {
	query := "UPDATE Cameras SET active=? WHERE id_camera=?"
	db := get()

	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(act, id)
	if err != nil {
		return 0, err
	}

	id_cam, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	defer stmt.Close()
	db.Close()
	return id_cam, nil
}

func UpdateTokenCon(idCam int, token_con string, user int) int64 {
	query := "UPDATE Cameras SET token_session_consumer=? WHERE id_camera=? and users_id=?"
	db := get()

	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	result, err := stmt.Exec(token_con, idCam, user)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	defer stmt.Close()
	db.Close()
	return id
}
func UpdateTokenCam(idCam int, user_id int, token_cam string) int64 {
	query := "UPDATE Cameras SET token_session_camera=? WHERE id_camera=? and users_id=?"
	db := get()

	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	result, err := stmt.Exec(token_cam, idCam, user_id)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	db.Close()
	defer stmt.Close()
	return id
}

func GetCamsByUser(user int) ([]m.Camera, error) {
	query := "SELECT * FROM Cameras WHERE users_id=?"
	db := get()
	cams := make([]m.Camera, 0)
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	result, err := stmt.Query(user)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var row m.Camera
		var date string
		//user, active, created, loc, url, token_session_camera, token_session_consumer, id_camera
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.Url, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			return nil, err
		}
		cams = append(cams, row)
	}
	defer db.Close()
	return cams, nil
}

func GetTokenCam(idCam int, idUser int) string {
	query := "SELECT * FROM Cameras WHERE id_camera=? AND users_id=?"
	db := get()
	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer stmt.Close()
	result, err := stmt.Query(idCam, idUser)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	var row m.Camera
	var date string

	for result.Next() {
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.Url, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	defer db.Close()
	fmt.Println(row.T_s_cam)
	return row.T_s_cam
}

func GetTokenCon(idCam int, idUser int) string {
	query := "SELECT * FROM Cameras WHERE id_camera=? AND users_id=?"
	db := get()
	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer stmt.Close()
	result, err := stmt.Query(idCam, idUser)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	var row m.Camera
	var date string

	for result.Next() {
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.Url, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	defer db.Close()
	fmt.Println(row.T_s_con)
	return row.T_s_con
}
func DeleteCam(idCam int, idUser int) (sql.Result, error) {
	query := "DELETE FROM Cameras WHERE id_camera=? and users_id=?"

	db := get()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Exec(idCam, idUser)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	return rows, nil

}

func GetNextCamIdByUser(currentUser int) (int, error) {
	query := "SELECT * FROM Cameras c WHERE users_id=? ORDER BY id_camera DESC LIMIT 1"

	db := get()

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Query(currentUser)

	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var row m.Camera
	var date string

	for result.Next() {
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.Url, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			return 0, err
		}
	}

	return row.Id_cam + 1, nil

}
