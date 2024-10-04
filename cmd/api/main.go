package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main(){
	ctx:= context.Background()
	if err:= run(ctx); err != nil{
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run(ctx context.Context) error{
	dbUser:= os.Getenv("DATABASE_USER")
	dbPassword:= os.Getenv("DATABASE_PASSWORD")
	dbHost:= os.Getenv("DATABASE_HOST")
	dbPort:= os.Getenv("DATABASE_PORT")
	dbName:= "coursesDB"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err:= sql.Open("postgres", dsn)
	if err != nil{
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer db.Close()
	err= db.Ping()
	if err != nil{
		log.Fatalf("Failed to open a DB connection on ping: %v", err)
	}else{
		fmt.Print("Success, connected to DB")
	}

	return nil
}