package obj

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/danieltmartin/ray-tracer/primitive"
	"github.com/danieltmartin/ray-tracer/tuple"
)

type Parser struct {
	rootGroup    *primitive.Group
	currentGroup *primitive.Group
	vertices     []tuple.Tuple
	normals      []tuple.Tuple
}

func newParser() Parser {
	group := primitive.NewGroup()
	return Parser{group, group, nil, nil}
}

func Parse(r io.Reader) (primitive.Primitive, error) {
	parsed, err := parse(r)
	return parsed.rootGroup, err
}

func parse(r io.Reader) (Parser, error) {
	file := newParser()
	sc := bufio.NewScanner(r)

	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		line, _, _ = strings.Cut(line, "#") // strip comments
		tokens := strings.Fields(line)
		err := file.parseLine(tokens)
		if err != nil {
			return file, fmt.Errorf("line %v: %v", lineNum, err)
		}
	}

	return file, nil
}

func (p *Parser) parseLine(tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	recType := tokens[0]
	var err error
	switch recType {
	case "v":
		err = p.parseVertex(tokens[1:])
	case "vn":
		err = p.parseNormal(tokens[1:])
	case "f":
		err = p.parseFace(tokens[1:])
	case "g":
		group := primitive.NewGroup()
		p.rootGroup.Add(group)
		p.currentGroup = group
	}
	return err
}

func (p *Parser) parseVertex(tokens []string) error {
	if len(tokens) < 3 {
		return fmt.Errorf("expected at least 3 arguments for vertex")
	}

	x, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}
	y, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}
	z, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}

	p.vertices = append(p.vertices, tuple.NewPoint(x, y, z))

	return nil
}

func (p *Parser) parseNormal(tokens []string) error {
	if len(tokens) < 3 {
		return fmt.Errorf("expected at least 3 arguments for normal")
	}

	x, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}
	y, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}
	z, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return fmt.Errorf("bad float: %v", err)
	}

	p.normals = append(p.normals, tuple.NewVector(x, y, z))

	return nil
}

type vertex struct {
	vertexIndex, normalIndex int
}

func (p *Parser) parseFace(tokens []string) error {
	if len(tokens) < 3 {
		return fmt.Errorf("expected at least 3 arguments for face")
	}

	var vertexReferences []vertex

	for _, token := range tokens {
		tokens = strings.Split(token, "/")
		if len(tokens) == 0 {
			return fmt.Errorf("expected at least a vertex reference for face")
		}
		v, err := strconv.ParseUint(tokens[0], 10, 64)
		if err != nil {
			return fmt.Errorf("bad vertex reference: %v", err)
		}
		if int(v) > len(p.vertices) {
			return fmt.Errorf("vertex does not exist")
		}

		var vertex vertex
		vertex.vertexIndex = int(v)
		vertex.normalIndex = -1

		if len(tokens) >= 3 {
			n, err := strconv.ParseUint(tokens[2], 10, 64)
			if err != nil {
				return fmt.Errorf("bad normal reference: %v", err)
			}
			if int(v) > len(p.vertices) {
				return fmt.Errorf("normal does not exist")
			}
			vertex.normalIndex = int(n)
		}

		vertexReferences = append(vertexReferences, vertex)
	}

	triangles := p.triangulate(vertexReferences)
	p.currentGroup.Add(triangles...)

	return nil
}

func (p *Parser) triangulate(vertexReferences []vertex) []primitive.Primitive {
	triangles := make([]primitive.Primitive, len(vertexReferences)-2)

	for i := 1; i < len(vertexReferences)-1; i++ {
		p1 := p.vertices[vertexReferences[0].vertexIndex-1]
		p2 := p.vertices[vertexReferences[i].vertexIndex-1]
		p3 := p.vertices[vertexReferences[i+1].vertexIndex-1]
		if vertexReferences[0].normalIndex != -1 {
			n1 := p.normals[vertexReferences[0].normalIndex-1]
			n2 := p.normals[vertexReferences[i].normalIndex-1]
			n3 := p.normals[vertexReferences[i+1].normalIndex-1]
			tri := primitive.NewSmoothTriangle(p1, p2, p3, n1, n2, n3)
			triangles[i-1] = &tri
		} else {
			tri := primitive.NewTriangle(p1, p2, p3)
			triangles[i-1] = &tri
		}
	}

	return triangles
}
