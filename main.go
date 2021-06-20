// GoParser project main.go
package main

func exdemo(a, b int) (int, int) {
	c := a
	a = b
	b = c
	return a, b
}

func main() {
	ParseGoSrc("demo.go.txt")
}
