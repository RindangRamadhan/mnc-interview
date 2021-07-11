package example

import "net/http"

type exampleHandlers struct{}

func (*exampleHandlers) List(w http.ResponseWriter, r *http.Request) {

}

var ExampleHandlers = &exampleHandlers{}
