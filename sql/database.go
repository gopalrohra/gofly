package sql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/gopalrohra/flyapi/env"
	"github.com/gopalrohra/flyapi/transformers"
	grpcdb "github.com/gopalrohra/grpcdb/grpc_database"
	"google.golang.org/grpc"
)

const (
	address = "localhost:3099"
)

type Database struct {
	dbInfo       *grpcdb.DatabaseInfo
	queryBuilder queryBuilder
}

func NewDatabase() Database {
	return Database{dbInfo: dbInfoFromEnv()}
}
func dbInfoFromEnv() *grpcdb.DatabaseInfo {
	return &grpcdb.DatabaseInfo{User: env.Config["DB_USER"], Password: env.Config["DB_PASSWORD"], Name: env.Config["DB_NAME"], Host: env.Config["DB_HOST"], Port: env.Config["DB_PORT"]}
}

func (db *Database) CreateDatabase() error {
	grpcClient := GRPCClient{}
	client, err := grpcClient.newClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, e := client.CreateDatabase(ctx, db.dbInfo)
	if e != nil {
		return e
	}
	if strings.ToLower(r.Status) != "success" {
		return errors.New("Something went wron.")
	}
	return nil
}
func (db *Database) Scan(i interface{}, queryClauses ...[]string) error {
	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer grpcCancel()
	conn, err := grpc.DialContext(grpcCtx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Couldn't connect to grpc server! Error: %v\n", err.Error())
		return err
	}
	defer conn.Close()
	client := grpcdb.NewGRPCDatabaseClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	binder := dataBinder{target: i}
	sq := db.queryBuilder.selectQuery(db.dbInfo, i, queryClauses...)
	return binder.bind(client.ExecuteSelect(ctx, sq))
}
func (db *Database) Insert(i interface{}) error {
	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer grpcCancel()
	conn, err := grpc.DialContext(grpcCtx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Couldn't connect to grpc server! Error: %v\n", err.Error())
		return err
	}
	defer conn.Close()
	client := grpcdb.NewGRPCDatabaseClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	iq := db.queryBuilder.insertQuery(db.dbInfo, i)
	iqr, err := client.ExecuteInsert(ctx, iq)
	if err != nil {
		return err
	}
	if strings.ToLower(iqr.Status) != "success" {
		return errors.New(("Something went wrong"))
	}
	fmt.Println(iqr.InsertedId)
	setReturnId(i, iqr.InsertedId)
	return nil
}
func setReturnId(i interface{}, insertedID string) {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		fmt.Println("Invalid target, must be ptr")
		return
	}
	if reflect.TypeOf(i).Elem().Kind() == reflect.Struct {
		v := reflect.ValueOf(i).Elem()
		f := v.FieldByName("ID")
		fmt.Printf("Is f: zero: %v\n", f.IsZero())
		if !f.IsZero() && f.CanSet() {
			transformers.Transformers[f.Type().String()](f, insertedID)
		}
	}
}

func (db *Database) Update(i interface{}, clauses []string) error {
	grpcClient := GRPCClient{}
	client, err := grpcClient.newClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uq := db.queryBuilder.updateQuery(*db.dbInfo, i, clauses)
	uqr, err := client.ExecuteUpdate(ctx, &uq)
	if err != nil {
		return err
	}
	if strings.ToLower(uqr.Status) != "success" {
		return errors.New("Something went wrong")
	}
	return nil
}
func (db *Database) CreateTable(tableName string, columns []string) error {
	grpcClient := GRPCClient{}
	client, err := grpcClient.newClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tableRequest := &grpcdb.TableRequest{Info: db.dbInfo, Name: tableName, ColumnDef: columns}
	tableResponse, err := client.CreateTable(ctx, tableRequest)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while creating table: %v", err))
		return err
	}
	if strings.ToLower(tableResponse.GetStatus()) != "success" {
		fmt.Println(tableResponse.GetDescription())
		return errors.New("Something went wrong")
	}
	fmt.Println(tableResponse.GetDescription())
	return nil
}

func (db *Database) AlterTable() {

}
