# Gotik
Kotik bot written in go

## Build and run on Windows

1. Download and run [Go installer](https://golang.org/dl/)
2. Download and run [MSYS2 installer](http://sourceforge.net/projects/msys2/)
    - Uncheck "Run MSYS2 32/64bit now"
3. Open the MSYS2 "MinGW-w64 Win32/64 Shell" from the start menu to install additional dependencies
    1. `pacman -Syu`
    2. `pacman -Syy git pkg-config mingw-w64-x86_64-toolchain mingw-w64-x86_64-opus mingw-w64-x86_64-ffmpeg`
4. Fetch repo, `git clone https://github.com/Mixaill/Gotik/ .`
5. Create a GOPATH (skip if you already have a GOPATH you want to use)
6. Configure environment for building Opus

````shell
export CGO_LDFLAGS="$(pkg-config --libs opus)"
export CGO_CFLAGS="$(pkg-config --cflags opus)"
````

7. Get go dependencies

`go get -tags nopkgconfig .`

8. Build Kotik

`go build -o Kotik.exe .`

9. Run Kotik

`./Kotik.exe`