import { z } from "zod";

/**
 * Thing corresponds to Go type examples.Thing (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * Thing is what we are interested in
 */
export const Thing = z.object({
    Name: z.string(),
    Count: z.number().int(),
});
export type Thing = z.infer<typeof Thing>;

/**
 * ExampleResult corresponds to Go type examples.ExampleResult (in package "github.com/softwaretechnik-berlin/goats/gotypes/examples").
 * The comment on the original Go type follows.
 *
 * ExampleResult is the result type for some call
 */
export const ExampleResult = z.object({
    Message: z.string(),
    Items: z.array(Thing).nullable().transform(a => a ?? []),
});
export type ExampleResult = z.infer<typeof ExampleResult>;
