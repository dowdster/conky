#!/bin/sh

cd ./cal

echo "Building Calendar..."
go build ./drawcal.go
go build ./calendar.go
go build ./myslice.go
cd ..
echo ""

echo "Building Mail..."
go build ./dateparse.go
go build ./gomail.go
cd ..
echo ""

echo "Building Tasks..."
go build ./tasks.go
cd ..

echo "Done"

