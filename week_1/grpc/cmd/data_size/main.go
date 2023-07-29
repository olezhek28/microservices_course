package main

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/proto"

	desc "github.com/olezhek28/microservices_course/week_1/grpc/pkg/note_v1"
)

func main() {
	session := &desc.NoteInfo{
		Title:    gofakeit.BeerName(),
		Content:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	dataJson, _ := json.Marshal(session)
	fmt.Printf("\n\ndataJson len %d byte \n%v\n", len(dataJson), dataJson)

	dataPb, _ := proto.Marshal(session)
	fmt.Printf("dataPb len %d byte \n%v\n", len(dataPb), dataPb)
}
