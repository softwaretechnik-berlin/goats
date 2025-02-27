// Package comments provides facilities for easily loading Go doc comments.
package comments

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/softwaretechnik-berlin/goats/gotypes/goinsp"
	//"github.com/softwaretechnik-berlin/goats/internal/domain"
)

type Loader interface {
	Load(t goinsp.Type) string
}

type loader struct {
	packagesConfig *packages.Config
	packages       map[goinsp.ImportPath]packageInfo
}

type packageInfo struct {
	typeInfos map[goinsp.TypeName]typeInfo
}

type typeInfo struct {
	declaration typeDeclaration
	methods     map[string]functionDeclaration
}

type typeDeclaration struct {
	genDecl  *ast.GenDecl
	typeSpec *ast.TypeSpec
	position token.Position
}

type functionDeclaration struct {
	funcDecl *ast.FuncDecl
	position token.Position
}

func (d typeDeclaration) CommentGroup() *ast.CommentGroup {
	if d.typeSpec.Doc != nil {
		return d.typeSpec.Doc
	}
	if len(d.genDecl.Specs) == 1 {
		return d.genDecl.Doc
	}
	return nil
}

func NewLoader() Loader {
	return loader{
		// For loading type comments: &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax},
		&packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | ^0, Tests: true},
		make(map[goinsp.ImportPath]packageInfo),
	}
}

func (l loader) Load(t goinsp.Type) string {
	info := l.typeInfo(t)
	commentGroup := info.declaration.CommentGroup()
	return commentGroup.Text()
}

func (l loader) LoadMethod(t goinsp.Type, name string) string {
	info := l.typeInfo(t)
	method, ok := info.methods[name]
	if !ok {
		panic(fmt.Sprintf("can't find %s for %s in %#v", name, t.Name(), info.methods))
	}
	commentGroup := method.funcDecl.Doc
	return commentGroup.Text()
}

func (l loader) typeInfo(t goinsp.Type) typeInfo {
	pkg, ok := l.packages[t.PkgPath()]
	if !ok {
		l.loadTypesFromPackage(t.PkgPath())
		pkg, ok = l.packages[t.PkgPath()]
		if !ok {
			panic(t)
		}
	}
	if td, ok := pkg.typeInfos[t.WithoutTypeArguments().Name()]; ok {
		return td
	}
	panic(t)
}

func (l loader) loadTypesFromPackage(path goinsp.ImportPath) {
	path = goinsp.ImportPath(strings.TrimSuffix(string(path), "_test"))
	pkgs, err := packages.Load(l.packagesConfig, string(path))
	if err != nil {
		panic(err)
	}
	if slices.IndexFunc(pkgs, func(p *packages.Package) bool { return p.PkgPath == string(path) }) == -1 {
		panic(pkgs)
	}
	for _, pkg := range pkgs {
		typeInfos := l.collectTypeInfo(pkg)
		l.packages[goinsp.ImportPath(pkg.PkgPath)] = packageInfo{typeInfos}
	}
}

func (l loader) collectTypeInfo(pkg *packages.Package) map[goinsp.TypeName]typeInfo {
	typeInfos := make(map[goinsp.TypeName]typeInfo)
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case (*ast.GenDecl):
				genDecl := decl
				for _, spec := range genDecl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						name := goinsp.TypeName(spec.Name.Name)
						info := typeInfos[name]
						info.declaration = typeDeclaration{genDecl, spec, pkg.Fset.Position(spec.Pos())}
						typeInfos[name] = info
					case *ast.ImportSpec, *ast.ValueSpec:
						// not currently interested in imports or values
					default:
						panic(spec)
					}
				}
			case (*ast.FuncDecl):
				if decl.Recv != nil {
					if decl.Recv.NumFields() != 1 {
						panic(decl)
					}
					receiver := decl.Recv.List[0].Type
					receiverName := nameOf(receiver)
					info := typeInfos[receiverName]
					if info.methods == nil {
						info.methods = make(map[string]functionDeclaration)
					}

					info.methods[decl.Name.Name] = functionDeclaration{decl, pkg.Fset.Position(decl.Pos())}
					typeInfos[receiverName] = info
				}
			default:
				panic(decl)
			}
		}
	}
	return typeInfos
}

func nameOf(expr ast.Expr) goinsp.TypeName {
	for {
		switch typed := expr.(type) {
		case *ast.Ident:
			return goinsp.TypeName(typed.Name)
		case *ast.StarExpr:
			expr = typed.X
		case *ast.IndexExpr:
			expr = typed.X
		case *ast.IndexListExpr:
			expr = typed.X
		default:
			panic(fmt.Sprintf("%v, i.e. %#v", expr, expr))
		}
	}
}
