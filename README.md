# Spaceship AAR
Using gomobile to generate the native AAR for Android. (or IOS?)  

## Usage
You may use it in your program, that's pretty convenient than executing the binary package.

## Generate
You have to ensure that you have gomobile already installed in your go mod, or using the following command to install it.
- for install gomobile in this project
```bash
go get golang.org/x/mobile/bind
```

Now, you can start binding this library.

- for multi platform
```bash
gomobile bind -target "android/arm64,android/amd64" -ldflags "-s -w" .
```
- for specified platform
```bash
gomobile bind -target android/arm64 -ldflags "-s -w" .
```
 
## Limitation
Since the gomobile only support very basic few types in class converter, it's better to pass a full configuration string
which represents the whole config structure in json format.

## Disclaimer
This repo or its generated library are for study purpose only, absolutely without any warranty, use it at your own risk.