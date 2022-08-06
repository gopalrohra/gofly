package sql

import (
	"context"
	"log"
	"time"

	grpcdb "github.com/gopalrohra/grpcdb/grpc_database"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	con *grpc.ClientConn
	err error
}

func (client *GRPCClient) newClient() (grpcdb.GRPCDatabaseClient, error) {
	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer grpcCancel()
	client.con, client.err = grpc.DialContext(grpcCtx, address, grpc.WithInsecure(), grpc.WithBlock())
	if client.err != nil {
		log.Printf("Couldn't connect to grpc server! Error: %v\n", client.err.Error())
		return nil, client.err
	}
	return grpcdb.NewGRPCDatabaseClient(client.con), nil
}
func (client *GRPCClient) close() {
	if client.err == nil {
		client.con.Close()
	}
}
