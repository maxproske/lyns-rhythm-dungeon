.PHONY: deps build run

deps:
	brew list sdl2{,_mixer,_image,_ttf} || brew install sdl2{,_mixer,_image,_ttf}
	go mod tidy

# -w -s reduces binary size by ~1.3MB
build:
	go build -ldflags "-w -s" -o lynsrd

run:
	go run main.go
