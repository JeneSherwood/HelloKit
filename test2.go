package main

type name interface {
}

func (receiver concatResponse) name() {

	// Implement the name method for concatResponse
	// ...
	// Your implementation here
}

func main() {
	var response concatResponse
	response.name()
}
