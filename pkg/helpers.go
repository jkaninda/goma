package pkg

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/jkaninda/goma/util"
	"log"
)

func Intro() {
	nameFigure := figure.NewFigure("Goma", "", true)
	nameFigure.Print()
	fmt.Printf("Version: %s\n", util.FullVersion())
	fmt.Println("Copyright (c) 2024 Jonas Kaninda")
	log.Println("Starting Goma Gateway...")
}
