package draw

type Path struct {
	Points []Point
}

func NewPath() Path {
	return Path{}
}

func (p Path) AppendPoint(point Point) Path {
	p.Points = append(p.Points, point)
	return p
}

func (p Path) RemovePoint(number int) Path {
	if number < 1 || number > len(p.Points) {
		return p
	}

	idx := number - 1
	p.Points = append(p.Points[:idx], p.Points[idx+1:]...)
	return p
}

func (p Path) Length() int {
	return len(p.Points)
}

func (p Path) GetPointNumber(number int) Point {
	if number < 1 || number > len(p.Points) {
		return Point{}
	}
	return p.Points[number-1]
}

func (p Path) Copy() Path {
	pathcopy := Path{}
	pathcopy.Points = []Point{}
	for _, p := range p.Points {
		pathcopy.Points = append(pathcopy.Points, p)
	}
	return pathcopy
}

func (p Path) Offset(offX, offY float64) Path {
	for i, pt := range p.Points {
		p.Points[i] = pt.Add(offX, offY)
	}
	return p
}

func (p Path) GetBoundingBox() BoundingBox {
	bbox := BoundingBox{}

	minX := 0.0
	maxX := 0.0
	minY := 0.0
	maxY := 0.0
	for idx, p := range p.Points {
		if idx == 0 {
			minX = p.X
			maxX = p.X
			minY = p.Y
			maxY = p.Y
			continue
		}

		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	bbox.X = minX
	bbox.Y = minY
	bbox.Width = maxX - minX
	bbox.Height = maxY - minY
	return bbox
}

type BoundingBox struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}
