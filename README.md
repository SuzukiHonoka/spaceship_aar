# Spaceship AAR
Using gomobile to generate the native AAR for Android. 

## Dependencies
- Android NDK
- gomobile

After you installed gomobile, you should use following command to initialize it:
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

## Generate
You have to ensure that you have gomobile already installed in your go mod, or using the following command to install it.

**Install gomobile in this project**
```bash
go get golang.org/x/mobile/bind
```

Now, you can start binding this library.

**Binding multi platform**
```bash
gomobile bind -androidapi 29 -target "android/arm64,android/amd64" -ldflags "-s -w" .
```

**Binding specified platform**
```bash
gomobile bind -androidapi 29 -target android/arm64 -ldflags "-s -w" .
```

## Usage
You may use it in your program, that's pretty convenient than executing the binary file.
 
## Limitation
Since the gomobile only support very basic few types in class converter, it's better to pass a full configuration string
which represents the whole config structure in json format.

## Disclaimer
This repo or its generated library are for study purpose only, absolutely without any warranty, use it at your own risk.