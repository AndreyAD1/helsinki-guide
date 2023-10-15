package main

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/translator"
)


func main() {
	ctx := context.Background()
	translator.Run(ctx)
}