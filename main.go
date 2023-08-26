package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/restuwahyu13/grpc-gateway/stubs/users"
)

const (
	TCP        = "tcp"
	UDPConn    = "udp"
	UNIXConn   = "unix"
	UNIXPACKET = "unixpacket"
)

type serverUserService struct {
	users.UnimplementedUsersServer
}

type clientUserService struct {
	client users.UsersClient
}

func main() {
	go GRPCServer()
	RestAPI()
}

/**
================================================================
= REST API TERITORY
================================================================
*/

func RestAPI() {
	var (
		server users.UnimplementedUsersServer = users.UnimplementedUsersServer{}
		client users.UsersClient              = GRPCClient()
		mux    *runtime.ServeMux              = runtime.NewServeMux()
	)

	err := users.RegisterUsersHandlerServer(context.Background(), mux, &server)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = users.RegisterUsersHandlerClient(context.Background(), mux, client)
	if err != nil {
		log.Fatal(err)
		return
	}

	mux.HandlePath(http.MethodPost, "/", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		var (
			req users.PingDTO
			res users.ApiResponse
		)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		rpc, err := client.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := rpc.SendMsg(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := rpc.RecvMsg(&res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(&res)
	})

	http.ListenAndServe(":3000", mux)
}

/**
================================================================
= RPC CLIENT TERITORY
================================================================
*/

func GRPCClient() users.UsersClient {
	ls, err := grpc.Dial(":30000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	client := users.NewUsersClient(ls)
	return &clientUserService{client: client}
}

func (h *clientUserService) Ping(ctx context.Context, opts ...grpc.CallOption) (users.Users_PingClient, error) {
	login, err := h.client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return login, nil
}

/**
================================================================
= RPC SERVER TERITORY
================================================================
*/

func GRPCServer() {
	ls, err := net.Listen(TCP, ":30000")
	if err != nil {
		log.Fatal(err)
		return
	}

	server := grpc.NewServer()
	users.RegisterUsersServer(server, &serverUserService{})
	server.Serve(ls)
}

func (h *serverUserService) Ping(stream users.Users_PingServer) error {
	var (
		req users.PingDTO
		res users.ApiResponse
	)

	if err := stream.RecvMsg(&req); err != nil {
		defer log.Fatal(err)
		return err
	}

	res.StatCode = http.StatusOK
	res.StatMessage = "Success"

	if err := stream.Send(&res); err != nil {
		defer log.Fatal(err)
		return err
	}

	fmt.Println("REQUEST FROM CLIENT: ", req.Test)

	return nil
}
