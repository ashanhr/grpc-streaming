/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"context"
	"log"
	"net"
	"time"

	greeter "github.com/wso2/choreo-samples/go-grpc/pkg"
	"google.golang.org/grpc"
)

type server struct {
	greeter.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeter.HelloRequest) (*greeter.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &greeter.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) StreamGreetings(in *greeter.HelloRequest, stream greeter.Greeter_StreamGreetingsServer) error {
	for i := 0; i < 10; i++ {
		if err := stream.Send(&greeter.HelloReply{Message: "Hello " + in.GetName()}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}