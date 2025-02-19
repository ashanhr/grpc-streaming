/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
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
	"io"
	"log"
	"os"
	"time"

	greeter "github.com/wso2/choreo-samples/go-grpc/pkg"
	"google.golang.org/grpc"
)

func main() {
	target := os.Getenv("GREETER_SERVICE")
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := greeter.NewGreeterClient(conn)
	name := "Choreo"

	// Unary call
	r, err := c.SayHello(context.Background(), &greeter.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	// Streaming call
	stream, err := c.StreamGreetings(context.Background(), &greeter.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not stream greetings: %v", err)
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive a greeting: %v", err)
		}
		log.Printf("Streaming Greeting: %s", reply.GetMessage())
		time.Sleep(1 * time.Second)
	}
}