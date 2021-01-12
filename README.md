# Visual Simulator Module

A trading terminal using a game engine. Using the power of graphics to quickly draw graphs and make decissions on how to
trade.

## Setup

### Install Go

Visit https://golang.org/dl/ and download the latest version for your OS.

### Run

#### Executable

Mac: `run.command`  
Windows: `run.bat`  
Linux: `run.sh`

#### Terminal

Make sure environment variables are set.

```
go run ./cmd/sim
```

## Environment

```
API_URL=https://localhost:5000 # url to market data
```