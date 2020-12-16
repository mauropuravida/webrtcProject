package factory

import (
	"log"
	/interfaces
	./model
)
//importar modelo y las interfaces

func FactoryDao(e string) interfaces.UserDao {
	var i interfaces.UserDao
	switch e {
	case "mysql":
		i = mysql.UserImplMysql{}
	default:
		log.Fatalf("El motor %s no esta implementado", e)
		return nil
	}

	return i
}
