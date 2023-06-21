package services

func GetDBCredentials() (string, string) {
	//dbType := os.Getenv("DB_TYPE")
	//dbName := os.Getenv("DB_NAME")
	//dbUser := os.Getenv("DB_USER")
	//dbPasswd := os.Getenv("DB_PASSWD")
	//dbHost := os.Getenv("DB_HOST")
	//dbPort := os.Getenv("DB_PORT")
	//return dbType, dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	return "mysql", "user:topor115ZG16P@tcp(localhost:3306)/talkersDB"
}
