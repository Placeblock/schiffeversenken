package main

import "schiffeversenken/ws"

func main() {

	ws.StartServer()

	<-make(chan int)
}
