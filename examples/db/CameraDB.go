package db
import ("time"
	"log"
	"net/http"
	"database/sql"
	 m "project/webrtcProject/examples/models"
)




func InsertCam(user int, loc string) (int64,error) {
	query := "INSERT INTO Cameras(user, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (?,?,?,?,?,?,?)"
	db := get()
	defer db.Close()
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
	defer db.Close()
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil,err
	}

	defer stmt.Close()
	result, err := stmt.Query(query,user)
	if err != nil {
		return nil,err
	}
	for result.Next() {
		var row m.Camera
		//user, active, created, loc, token_session_camera, token_session_consumer, id_camera
		err := result.Scan(&row.ID, &row.Active, &row.Loc, &row.T_s_cam, &row.T_s_con, &row.Id_cam)
		if err != nil {
			return nil, err
		}

		cams = append(cams, row)
	}
	return cams,nil
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
func deleteCam(w http.ResponseWriter, r *http.Request) (sql.Result ,error) {
	query := "DELETE FROM Cameras WHERE user=?"
	
	db := get()
	defer db.Close()

	user:= r.URL.Query().Get("id")

	log.Println("delete"+user)
	stmt, err := db.Prepare(query)
	if err != nil {
		return  nil, err
	}

	defer stmt.Close()

	rows, err :=stmt.Exec(user)

	if err != nil {
		return nil, err
	}
	
	return rows,nil


}
