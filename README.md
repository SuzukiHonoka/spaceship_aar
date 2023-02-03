# Spaceship AAR
Using gomobile to generate the native AAR for Android. (or IOS?)  

## Usage
You may use it in your Android program, that's pretty convenient than executing the program binary.

## Generate
You have to ensure that you have gomobile already installed in your go mod, then execute the following:

```bash
gomobile bind -target android/arm64 .
```
 
## Limitation
Since the gomobile only support very basic few types in class converter, it's better to pass a full configuration string
which represents the whole config structure in json format.

## Disclaimer
This repo or its generated library are for study purpose only, absolutely without any warranty, use it at your own risk.