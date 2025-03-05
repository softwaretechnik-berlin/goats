package gozod

import (
	"slices"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/parsing/comments"
	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp/reflective"
	"github.com/softwaretechnik-berlin/goats/gotypes/ts"
	"github.com/softwaretechnik-berlin/goats/gotypes/util"
	"github.com/softwaretechnik-berlin/goats/gotypes/zod"
)

type config struct {
	names                 map[goinsp.GenType]ts.Identifier
	unnamedTypes          map[goinsp.GenType]struct{}
	schemas               map[goinsp.GenType]func(resolver Resolver[Ref, zod.ZodType]) zod.ZodType
	templates             map[goinsp.GenType]string
	undiscriminatedUnions map[goinsp.GenType][]goinsp.Type
	discriminators        map[goinsp.GenType]JSONDiscriminator
	discriminatedUnions   map[goinsp.GenType]JSONDiscriminatedUnion
	transforms            map[goinsp.GenType]func(resolver Resolver[Ref, zod.ZodType]) ts.Source
	commentsLoader        comments.Loader
}

type JSONDiscriminator struct {
	Property string
	Value    string
}

type JSONDiscriminatedUnion struct {
	DiscriminatorProperty string
	Types                 []goinsp.Type
}

func newConfig(options ...Option) config {
	c := config{}
	for _, o := range options {
		o.apply(&c)
	}
	if c.commentsLoader == nil {
		c.commentsLoader = comments.NewLoader()
	}
	return c
}

type Option interface {
	apply(*config)
}

var _ Option = funcOption(nil)
var _ Option = TypeOptions{}

func WithName(t goinsp.GenType, name string) Option {
	return funcOption(func(c *config) {
		if c.names == nil {
			c.names = make(map[goinsp.GenType]ts.Identifier)
		}
		c.names[t] = ts.Identifier(name)
	})
}

func WithUnnamedType(t goinsp.GenType) Option {
	return funcOption(func(c *config) {
		if c.unnamedTypes == nil {
			c.unnamedTypes = make(map[goinsp.GenType]struct{})
		}
		c.unnamedTypes[t] = struct{}{}
	})
}

func WithResolvingSchema(t goinsp.GenType, schema func(resolver Resolver[Ref, zod.ZodType]) zod.ZodType) Option {
	return funcOption(func(c *config) {
		if c.schemas == nil {
			c.schemas = make(map[goinsp.GenType]func(resolver Resolver[Ref, zod.ZodType]) zod.ZodType)
		}
		c.schemas[t] = schema
	})
}

func WithSchema(t goinsp.GenType, schema zod.ZodType) Option {
	return WithResolvingSchema(t, func(_ Resolver[Ref, zod.ZodType]) zod.ZodType { return schema })
}

func WithTemplate(t goinsp.GenType, template string) Option {
	return funcOption(func(c *config) {
		if c.templates == nil {
			c.templates = make(map[goinsp.GenType]string)
		}
		c.templates[t] = template
	})
}

func WithUndiscriminatedUnion(t goinsp.GenType, types ...goinsp.Type) Option {
	return funcOption(func(c *config) {
		if c.undiscriminatedUnions == nil {
			c.undiscriminatedUnions = make(map[goinsp.GenType][]goinsp.Type)
		}
		c.undiscriminatedUnions[t] = types
	})
}

func WithDiscriminator(t goinsp.GenType, property, value string) Option {
	return funcOption(func(c *config) {
		if c.discriminators == nil {
			c.discriminators = make(map[goinsp.GenType]JSONDiscriminator)
		}
		c.discriminators[t] = JSONDiscriminator{property, value}
	})
}

func WithDiscriminatedUnion(t goinsp.GenType, discriminatorProperty string, types ...goinsp.Type) func(*config) {
	return func(c *config) {
		if c.discriminatedUnions == nil {
			c.discriminatedUnions = make(map[goinsp.GenType]JSONDiscriminatedUnion)
		}
		c.discriminatedUnions[t] = JSONDiscriminatedUnion{discriminatorProperty, types}
	}
}

func WithResolvingTransform(t goinsp.GenType, expr func(resolver Resolver[Ref, zod.ZodType]) ts.Source) Option {
	return funcOption(func(c *config) {
		if c.transforms == nil {
			c.transforms = make(map[goinsp.GenType]func(resolver Resolver[Ref, zod.ZodType]) ts.Source)
		}
		c.transforms[t] = expr
	})
}

func WithTransform(t goinsp.GenType, expr ts.Source) Option {
	return WithResolvingTransform(t, func(_ Resolver[Ref, zod.ZodType]) ts.Source { return expr })
}

func WithCommentsLoader(loader comments.Loader) Option {
	return funcOption(func(c *config) {
		c.commentsLoader = loader
	})
}

type TypeOptions struct {
	t       goinsp.GenType
	options []Option
}

func (o TypeOptions) apply(c *config) {
	for _, option := range o.options {
		option.apply(c)
	}
}

func (o TypeOptions) Unnamed() TypeOptions {
	return o.add(WithUnnamedType(o.t))
}

func (o TypeOptions) Schema(schema zod.ZodType) TypeOptions {
	return o.add(WithSchema(o.t, schema))
}

func (o TypeOptions) ResolvingSchema(schema func(resolver Resolver[Ref, zod.ZodType]) zod.ZodType) TypeOptions {
	return o.add(WithResolvingSchema(o.t, schema))
}

func (o TypeOptions) Template(template string) TypeOptions {
	return o.add(WithTemplate(o.t, template))
}

func (o TypeOptions) add(options ...Option) TypeOptions {
	return TypeOptions{o.t, append(slices.Clip(o.options), options...)}
}

func (o TypeOptions) Named(name string) TypeOptions {
	return o.add(WithName(o.t, name))
}

func (o TypeOptions) ResolvingTransform(f func(resolver Resolver[Ref, zod.ZodType]) ts.Source) TypeOptions {
	return o.add(WithResolvingTransform(o.t, f))
}

func (o TypeOptions) Transform(f ts.Source) TypeOptions {
	return o.add(WithTransform(o.t, f))
}

func (o TypeOptions) Transformf(format string, as ...any) TypeOptions {
	return o.add(WithResolvingTransform(o.t, func(resolver Resolver[Ref, zod.ZodType]) ts.Source {
		return ts.Sourcef(format, util.Map(as, func(value any) ts.Source {
			switch value := value.(type) {
			case ts.Source:
				return value
			case goinsp.Type:
				return resolver.Resolve(PlainRef(value)).TypeScript()
			case Ref:
				return resolver.Resolve(value).TypeScript()
			default:
				panic(value)
			}
		})...)
	}))
}

func (o TypeOptions) UndiscriminatedUnionOf(disjuncts ...goinsp.Type) Option {
	return o.add(WithUndiscriminatedUnion(o.t, disjuncts...))
}

func ForType(t goinsp.GenType) TypeOptions {
	return TypeOptions{t, nil}
}

func When[T any]() TypeOptions {
	return ForType(reflective.TypeFor[T]())
}

func WhenGeneric[T any]() TypeOptions {
	t := reflective.TypeFor[T]()
	if generic := t.WithoutTypeArguments(); generic != t {
		return ForType(generic)
	}
	panic(t)
}

type funcOption func(c *config)

func (f funcOption) apply(c *config) {
	f(c)
}

func lookupConfig[T any](m map[goinsp.GenType]T, t goinsp.Type) (value T, ok bool) {
	value, ok = m[t]
	if !ok {
		value, ok = m[t.WithoutTypeArguments()]
	}
	return
}
