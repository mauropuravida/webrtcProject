package db
import ("time"
	
	"database/sql"
	 m "project/webrtcProject/examples/models"
	 "fmt"
)




func InsertCam(user int, loc string) (int64,error) {
	query := "INSERT INTO Cameras(user, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?)"
	db := get()
	
	stmt, err := db.Prepare(query)

	if err != nil {
		return 0,err
	}

	defer stmt.Close()
	result, err := stmt.Exec(user, false, time.Now() , loc , " ", " " , 0)
	if err != nil {
		return 0,err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0,err
	}

	return id,nil
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
		//user, active, created, loc, token_session_camera, token_session_consumer, id_camera
		err := result.Scan(&row.ID, &row.Active, &date, &row.Loc, &row.T_s_cam, &row.T_s_con, &row.Id_cam, &row.User_id)
		if err != nil {
			return nil, err
		}
		fmt.Println(row.ID)
		cams = append(cams, row)
	}
	defer db.Close()
	return cams,nil
}

func DeleteCam(idCam int, idUser int ) (sql.Result ,error) {
	query := "DELETE FROM Cameras WHERE user_id=? AND id=?"
	
	db := get()


	stmt, err := db.Prepare(query)
	if err != nil {
		return  nil, err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(idUser,idCam)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	return rows,nil


}
