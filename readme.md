
# Intel 8080 Emulator and Space Invaders arcade implementation

This project is trying to emulate i8080 cpu and also based on cpu implementation make something more impressive than executing instructions, for example, make ancient space invaders code written in assembly running.


![Alt Text](https://media0.giphy.com/media/v1.Y2lkPTc5MGI3NjExOXh1MWRrYXp6Zng2bmJyeHQwd29nd2FreDcwemtxdjIwMTVsZGZmciZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/ZGbnMqUNjuRueOioWm/giphy.gif)


## Requirements

[Go v1.13+](https://go.dev/dl/)
[Go SDL](https://github.com/veandco/go-sdl2?tab=readme-ov-file#requirements)



## Build
```bash
  go run build .
```

## How to run

You can either play the game or use debugger. 

To run debugger you have to specify path and debugger flag.

To run the game you have to simply specify -p flag
####  run emulator
```bash
  ./cpu-emulator
```
## Arguments

| Flag             | Description|
| ----------------- | ------------------------------------------------------------------ |
| -p | run space invaders |
| -r  | path to ROM |
| -d | run debugger |

## Example 
To run debugger type in terminal
```bash
  ./cpu-emulator -r [path to rom] -d 
```

# Key bindings
| Key             | Action description|
| ----------------- | ------------------------------------------------------------------ |
| A | move left |
| D  | move right |
| Space | shoot |
| W| Start|
|S | Insert coin|





