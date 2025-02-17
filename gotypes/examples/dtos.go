package examples

// Example1 is the result type for some call
type Example1 struct {
	// Message is the message the server produced for the request
	Message string
	// Items are the Things that we are interested in.
	Items []ChildThing1
}

// ChildThing1 is what we are interested in
type ChildThing1 struct {
	// Name is used to commonly refer to the thing
	Name string
	// Count multiple things.
	Count int
}

// Example2 is the result type for some call
type Example2 struct {
	// Message is the message the server produced for the request
	Message string `json:"message"`
	// Items are the Things that we are interested in.
	Items []ChildThing2 `json:"items"`
}

// ChildThing2 is what we are interested in
type ChildThing2 struct {
	// Name is used to commonly refer to the thing
	Name string `json:"name"`
	// Count multiple things.
	Count int `json:"count"`
}

// Example3 a struct containing a map
type Example3 struct {
	// Elements
	Elements map[string]int
}
