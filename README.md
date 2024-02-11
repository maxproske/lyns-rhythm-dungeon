# Lyn's Rhythm Dungeon

What's in a name?

## Why Go?

- Simple, readable, performant code.
- Compiles to native executables. Simplifies sharing games, and doesn't rely on just-in-time compilation.
- Garbage collection is very low latency, compared to Java or C# which focus on throughput.
- First-class support for concurrency, memory allocation, and unit testing.

## Getting Started (MacOS)

```sh
# Install Go
brew install go

# Install SDL2
brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
go get -v github.com/veandco/go-sdl2/{sdl,img,mix,ttf}

# Build
go build -o lynsrd
./lynsrd
```
