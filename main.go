package main

func main() {
	
	go Publisher()
	go Subscriber()
	select {}
}