package main

import (
	"math/rand"
	"time"

	"github.com/JosephSalisbury/vm/cmd"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	cmd.Execute()
}
