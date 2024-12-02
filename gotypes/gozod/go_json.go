package gozod

import (
	"strings"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
)

func forEachTopLevelJSONFieldAndEmbeddedType(
	t goinsp.Type,
	visitPropertyField func(name string, field goinsp.StructField, tag string),
	visitEmbeddedJSONType func(t goinsp.Type),
) {
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == `-` {
			continue
		}
		if field.Anonymous && tag == "" {
			visitEmbeddedJSONType(field.Type())
		} else if field.IsExported() {
			visitPropertyField(jsonPropertyName(field, tag), field, tag)
		}
	}
	return
}

func jsonPropertyName(field goinsp.StructField, tag string) string {
	name, _, _ := strings.Cut(tag, ",")
	if name != "" {
		return name
	}
	return field.Name
}

//func forEachJsonProperty(t goinsp.Type, visit func(name string, tag string, field goinsp.StructField)) goinsp.Type {
//	type visibility struct {
//		level      int
//		jsonTagged bool
//		count      int
//	}
//	visibilityByPropertyName := make(map[string]visibility)
//	forEachPotentialJSONProperty(0, t, func(nestingLevel int, field goinsp.StructField, jsonTag string, propertyName string) {
//		vis, ok := visibilityByPropertyName[propertyName]
//		jsonTagged := jsonTag == ""
//		if !ok || nestingLevel < vis.level || nestingLevel == vis.level && jsonTagged && !vis.jsonTagged {
//			visibilityByPropertyName[propertyName] = visibility{nestingLevel, jsonTagged, 1}
//		} else if vis.level == nestingLevel && vis.jsonTagged == jsonTagged {
//			visibilityByPropertyName[propertyName] = visibility{nestingLevel, jsonTagged, vis.count + 1}
//		}
//	})
//	propertyCount := 0
//	representation := forEachPotentialJSONProperty(0, t, func(nestingLevel int, field goinsp.StructField, tag string, name string) {
//		vis := visibilityByPropertyName[name]
//		if vis.level == nestingLevel && vis.jsonTagged == (tag == "") && vis.count == 1 {
//			propertyCount++
//			visit(name, tag, field)
//		}
//	})
//	if propertyCount == 0 {
//		return representation
//	}
//	return nil
//}

//func forEachPotentialJSONProperty(nestingLevel int, t goinsp.Type, visit func(nestingLevel int, field goinsp.StructField, tag string, name string)) (representationType goinsp.Type) {
//	for i := range t.NumField() {
//		field := t.Field(i)
//		tag := field.Tag.Get("json")
//		if tag == "-" {
//			continue
//		}
//		if field.Anonymous && tag == "" {
//			t := field.Type()
//			switch t.Kind() {
//			case reflect.Array:
//				representationType = t
//			case reflect.Struct:
//				representationType = forEachPotentialJSONProperty(nestingLevel+1, t, visit)
//			default:
//				panic(t)
//			}
//		} else if field.IsExported() {
//			visit(nestingLevel, field, tag, jsonPropertyName(field, tag))
//		}
//	}
//	return
//}
