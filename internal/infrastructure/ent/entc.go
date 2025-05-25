//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	err := entc.Generate("./schema", &gen.Config{
		Target:   "./entgen",
		Package:  "github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen",
		Features: gen.AllFeatures,
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
