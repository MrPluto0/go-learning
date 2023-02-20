package main

import (
	"log"
	"protoc/entity"

	"google.golang.org/protobuf/proto"
)

func main() {
	user := &entity.User{
		Name:   "gypsophlia",
		Avatar: "spike.png",
	}

	data, err := proto.Marshal(user)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	user2 := &entity.User{}
	err = proto.Unmarshal(data, user2)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	log.Printf("Name: %v %v", user.GetName(), user.Name)
}
