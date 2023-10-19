package main

import (
	"os"

	"github.com/adelblande/codepix/application/grpc"
	"github.com/adelblande/codepix/infrastructure/db"
	"github.com/jinzhu/gorm"
)

var database *gorm.DB

func main() {
	database = db.ConnectDb(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}