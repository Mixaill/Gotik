# Gotik
Kotik bot written in go


##Compile

##Build
### Windows

1. Install dependencies
    1. Base dependencies
        1. Download and run [Go installer](https://golang.org/dl/)
        2. Download and run [MSYS2 installer](http://sourceforge.net/projects/msys2/)
            - Uncheck "Run MSYS2 32/64bit now"
    2. Open the MSYS2 "MinGW-w64 Win32/64 Shell" from the start menu to install additional dependencies
        - `pacman -Syu`
        - `pacman -Syy mingw-w64-i686-toolchain git mingw-w64-i686-opus pkg-config mingw-w64-i686-ffmpeg`
2. Move to target directory
3. Fetch repo
    - `git clone https://github.com/Mixaill/Gotik/ .`
4. Create a GOPATH (skip if you already have a GOPATH you want to use)
    - `mkdir -p ./GOPATH`
    - `export GOPATH=$(pwd)/GOPATH`
5. Configure environment for building Opus
    - `export CGO_LDFLAGS="$(pkg-config --libs opus)"`
    - `export CGO_CFLAGS="$(pkg-config --cflags opus)"`
6. Get go dependencies
    - `go get -tags nopkgconfig .`
7. Build Kotik
    - `go build -o Kotik.exe .`
8. Run Kotik
    - `./Kotik.exe`
