package db
import ("time"
	
	"database/sql"
	 m "project/webrtcProject/examples/models"
	 "fmt"
)




func InsertCam(user int, loc string, url string) (int64,error) {
	
	query := "INSERT INTO Cameras(users_id, active, created, loc, url, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?,?)"
	fmt.Println("INSERTANDO")
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
	fmt.Println("error1")
	fmt.Println(err)
		return 0,err
	}

	defer stmt.Close()
	result, err := stmt.Exec(user, false, time.Now() , loc , url, " ", " " , 0)
	if err != nil {
	fmt.Println("error2")
	fmt.Println(err)
		return 0,err
	}

	id, err := result.LastInsertId()
	if err != nil {
	fmt.Println("error3")
	fmt.Println(err)
		return 0,err
	}

	return id,nil
}

func UpdateCam(idCam int, user int, loc string, url string) (int64,error) {
	query := "UPDATE Cameras SET users_id=?, loc=?, url=? WHERE id_camera=?"
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
	fmt.Println("error1")
	fmt.Println(err)	
		return 0,err
	}

	result, err := stmt.Exec(user, loc, url ,idCam)
	if err != nil {
	fmt.Println("error2")
	fmt.Println(err)	
		return 0,err
	}

	id, err := result.LastInsertId()
	if err != nil {
	fmt.Println("error3")
	fmt.Println(err)	
		return 0,err
	}
	
	defer stmt.Close()
	return id,nil
}


func UpdateActiveCam(act bool, id int) (int64,error) {
	query := "UPDATE Cameras SET active=? WHERE id_camera=?"
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
	fmt.Println("error1")
	fmt.Println(err)	
		return 0,err
	}

	result, err := stmt.Exec(id, act)
	if err != nil {
	fmt.Println("error2")
	fmt.Println(err)	
		return 0,err
	}

	id_cam, err := result.LastInsertId()
	if err != nil {
	fmt.Println("error3")
	fmt.Println(err)	
		return 0,err
	}
	
	defer stmt.Close()
	return id_cam,nil
}



func UpdateTokenCon(idCam int,  token_con string) (int64)  {
	query := "UPDATE Cameras SET token_session_connection=? WHERE id_camera=?"
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
	fmt.Println(err)	
		return 0
	}

	result, err := stmt.Exec(token_con, idCam)
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
	return id
}
func UpdateTokenCam(idCam int,  token_cam string) (int64) {
	query := "UPDATE Cameras SET token_session_camera=? WHERE id_camera=?"
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)	
		return 0
	}

	result, err := stmt.Exec(token_cam, idCam)
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
	return id
}


	
func GetCamsByUser(user int) ([]m.Camera,error) {
	query := "select * from cameras where users_id=?"
	db := get()
	cams := make([]m.Camera, 0)
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil,err
	}

	defer stmt.Close()
	result, err := stmt.Query(user)
	if err != nil {
		return nil,err
	}
	for result.Next() {
		var row m.Camera
		var date string
		//user, active, created, loc, url, token_session_camera, token_session_consumer, id_camera
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.Url, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			return nil, err
		}
		fmt.Println(row.ID)
		cams = append(cams, row)
	}
	defer db.Close()
	return cams,nil
}




	
func GetTokenCam(idcam int) (string) {
	query := "select token_session_camera from cameras where id=?"
	db := get()
	stmt, err := db.Prepare(query)

	if err != nil {
		return ""
	}

	defer stmt.Close()
	result, err := stmt.Query(idcam)

	if err != nil {
		return ""
	}
	
	var token string
    err = result.Scan(&token)
	defer db.Close()
	return token
}

func DeleteCam(idCam int) (sql.Result ,error) {
	query := "DELETE FROM Cameras WHERE  id=?"
	
	db := get()


	stmt, err := db.Prepare(query)
	if err != nil {
		return  nil, err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(idCam)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	return rows,nil


}
