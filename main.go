package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arsenypoga/canvas-api/pkg/api"
)

func main() {
	token := os.Getenv("CANVAS_OAUTH")
	fmt.Println(token)

	client := api.NewClient("nku", token)
	resp, err := client.GetTodo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}
