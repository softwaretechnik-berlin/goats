# goats

goats takes golang struct definitions and produces zod signatures to be used to 
parse and type the responses in typescript client code. tsgen is written as a library
that you can run and configure by invoking functions.

Here is a basic example consider these DTO structs:

~~~go
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
~~~

Now running goats from a main function in your code base. Configuring it in 
a typesafe way using code:

~~~golang
import (
    "github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
    "github.com/softwaretechnik-berlin/goats/gotypes/gozod"
)

mapper := gozod.NewMapper()

mapper.ResolveAll(  // ResolveAll will automatically resolve the `ChildThing1` by walking the dependency tree.
    reflective.TypeFor[Example1](),  
)

gozod.Generate(mapper, "basic_example.ts")
~~~

This will generate the following file with types that correspond to what it is serialised to. 

~~~typescript
import { z } from "zod";

/**
 * ChildThing1 corresponds to Go type examples.ChildThing1 (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * ChildThing1 is what we are interested in
 */
export const ChildThing1 = z.object({
    Name: z.string(),
    Count: z.number().int(),
});
export type ChildThing1 = z.infer<typeof ChildThing1>;

/**
 * Example1 corresponds to Go type examples.Example1 (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * Example1 is the result type for some call
 */
export const Example1 = z.object({
    Message: z.string(),
    Items: z.array(ChildThing1).nullable().transform(a => a ?? []),
});
export type Example1 = z.infer<typeof Example1>;
~~~

## json tags

tsgen also takes `json` tags into account, so this slightly modified example:

~~~go
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
~~~

will yield the following zod schema definition:

~~~typescript
import { z } from "zod";

/**
 * ChildThing2 corresponds to Go type examples.ChildThing2 (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * ChildThing2 is what we are interested in
 */
export const ChildThing2 = z.object({
    name: z.string(),
    count: z.number().int(),
});
export type ChildThing2 = z.infer<typeof ChildThing2>;

/**
 * Example2 corresponds to Go type examples.Example2 (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * Example2 is the result type for some call
 */
export const Example2 = z.object({
    message: z.string(),
    items: z.array(ChildThing2).nullable().transform(a => a ?? []),
});
export type Example2 = z.infer<typeof Example2>;
~~~

## Using Maps

Maps get mapped to the TypeScript record type: 

~~~go
// Example3 a struct containing a map
type Example3 struct {
	// Elements
	Elements map[string]int
}
~~~

This will yield the following zod type:

~~~typescript
import {z} from "zod";

/**
 * Example3 corresponds to Go type examples.Example3 (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * Example3 a struct containing a map
 */
export const Example3 = z.object({
    Elements: z.record(z.string(), z.number().int()).nullable().transform(r => r ?? {})
});

export type Example3 = z.infer<typeof Example3>;

~~~