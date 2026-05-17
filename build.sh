#! /bin/bash
cd traceinspector
go build -o ../front_src/bin/traceinspector.o ./cmd/traceinspector/main.go
cd ..