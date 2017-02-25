package main

import ("fmt" 
		"encoding/json"
		"net/http"
		"math"
		"errors")

type Coord struct {
	X int `json:"x"`
	Y int  `json:"y"`
	Z int  `json:"z"`
}

type Translation struct {
	X int `json:"x"`
	Y int  `json:"y"`
	Z int  `json:"z"`
}

type OrientedBlock struct {
	BlockID int	 `json:"blockID"`	// holds the ID of the block
	Parts []Coord `json:"parts"`	// holds a set of base coordinates of the pieces
	localMovements []Translation    // list of local movements that would still keep a homed block within the 4x4 bounding box
}

type Block struct {
	BlockID int
	Variations [24]*OrientedBlock
}

type Axis int
const (
	XAxis Axis = iota
	YAxis
	ZAxis)

const SpaceDimension int = 30


type Space struct {
	size int
	grid [][][]*OrientedBlock
}


// Some global variables. Forgive me

var blocks [6]*Block;

// Each block is described with a base orientation. We will have to calculate the 24 possible orientations

var block0 = OrientedBlock{BlockID:0,
	Parts:[]Coord{Coord{0,0,0}, Coord{1,0,0}, Coord{0,1,0},Coord{0,2,0}, Coord{0,3,0}, Coord{1,3,0}}}

var block1 = OrientedBlock{BlockID:1,
	Parts:[]Coord{Coord{1,0,0}, Coord{2,0,0}, Coord{3,0,0}, Coord{1,1,0}, Coord{1,2,0}, Coord{1,3,0},
		Coord{2,3,0}, Coord{1,0,1}, Coord{0,0,1}, Coord{0,0,2}, Coord{0,0,3}, Coord{0,1,1}}}

var block2 = OrientedBlock{BlockID:1,
	Parts:[]Coord{Coord{1,0,0}, Coord{2,0,0}, Coord{3,0,0}, Coord{0,1,0}, Coord{1,1,0}, Coord{3,1,0}, 
		Coord{0,2,0}, Coord{0,2,1}, Coord{3,2,0}, Coord{0,3,0}}}

var block3 = OrientedBlock{BlockID:3,
	Parts:[]Coord{Coord{0,0,0}, Coord{1,0,0}, Coord{2,0,0}, Coord{0,0,1}, Coord{0,1,1}, Coord{0,2,1}, Coord{0,3,1}, 
		Coord{2,1,0}, Coord{3,1,0}}}

var block4 = OrientedBlock{BlockID:4,
	Parts:[]Coord{Coord{0,0,0}, Coord{1,0,0}, Coord{2,0,0}, Coord{0,1,0}, Coord{2,1,0}, Coord{3,1,0}, 
			Coord{0,2,0}, Coord{3,2,0}, Coord{3,3,0}, Coord{0,1,1}, Coord{0,0,2}, Coord{0,1,2},
			Coord{0,1,3}, Coord{0,2,3}, Coord{0,3,3}}}

var block5 = OrientedBlock{BlockID:5,
	Parts:[]Coord{Coord{0,0,0}, Coord{1,0,0}, Coord{2,0,0}, Coord{0,1,0}, Coord{0,2,0},
			Coord{2,0,1}, Coord{0,1,1}, Coord{2,1,1}, Coord{2,2,1}, Coord{3,2,1}, Coord{2,2,2}}}


// every possible starting position for a block if we want to build the final cube in the positions 0,0,0 -> 3,3,3
var startingPositions = []Coord{Coord{-8,0,0},Coord{8,0,0}, Coord{0,-8,0}, Coord{0,8,0}, Coord{0,0,-8}, Coord{0,0,8}}



// gets the element at the x,y,z value
func (s* Space) GetElem(x int, y int, z int) (*OrientedBlock, error) {
	// move into positive numbers
	x = x + s.size/2
	y = y + s.size/2
	z = z + s.size/2
	if (x >= s.size) || (y >= s.size) || (z >= s.size) {
		error := errors.New("(Space) Get Elem, x,y, or z out of bounds")
		return nil, error
	}
	return s.grid[x][y][z], nil

}

// sets the element at the x,y,z value, returning error if out of bounds
// also returns the old value
func (s* Space) SetElem(x int, y int, z int, block *OrientedBlock) (*OrientedBlock, error) {
	x = x + s.size/2
	y = y + s.size/2
	z = z + s.size/2
	if (x >= s.size) || (y >= s.size) || (z >= s.size) {
		error := errors.New("(Space) Set Elem, x,y, or z out of bounds")
		return nil, error
	}
	oldValue := s.grid[x][y][z]
	s.grid[x][y][z] = block
	return oldValue, nil
}

func NewSpace(size int) (*Space) {
	var mySpace *Space = new(Space)
	mySpace.size = size
	mySpace.grid = make([][][]*OrientedBlock, size)
	for i:=0; i< size; i++ {
		mySpace.grid[i] = make([][] *OrientedBlock, size)
		for j:=0; j < size; j++ {
			mySpace.grid[i][j] = make([]*OrientedBlock, size)
			for k:=0; k < size; k++ {
				mySpace.grid[i][j][k] = nil   // probably not needed but not sure
			}
		}
	}

	return mySpace

}

// Attempts to insert every part of the block and returns true if it can and false if it can
// if it returns false, it writes nothing (so it does this in two passses)
func (s* Space) InsertBlock(block *OrientedBlock) (bool, error) {

	if (block == nil) {
		return false, errors.New("(Space ) received empty block")
	}

	for i:=0; i < len(block.Parts); i++ {
		elem, e := s.GetElem(block.Parts[i].X, block.Parts[i].Y, block.Parts[i].Z)
		if (e != nil) {
			return false, e
		}
		if (elem != nil) {
			return false, nil // there is a conflict
		}
	}
	// at this point, we can insert the block
	for i:=0; i < len(block.Parts); i++ {
		elem, e := s.SetElem(block.Parts[i].X, block.Parts[i].Y, block.Parts[i].Z, block)
		if (elem != nil) {
			return true, errors.New(fmt.Sprintf("(Space) Internal error when inserting, %s",e.Error()))
		}
	}
	return true, nil
}


// Given an oriented Block, return a Block with all 24 variations
func CreateBlockOrientations(baseBlock *OrientedBlock) (*Block) {
	b := new(Block)
	var v int = 0
	b.BlockID = baseBlock.BlockID
	b.Variations[0],v = baseBlock, v + 1

	var o *OrientedBlock = baseBlock
	b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 90)
	b.Variations[v].TranslateToQuadrantOne()
	v++
	b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 180)
	b.Variations[v].TranslateToQuadrantOne()
	v++
	b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 270)
	b.Variations[v].TranslateToQuadrantOne()
	v++

	for alpha := 90; alpha < 360; alpha+=90 {
		// should call for 90, 180, 270  (12 in total)
		o, _ = baseBlock.RotateAroundAxis(XAxis, alpha)
		b.Variations[v] = o
		b.Variations[v].TranslateToQuadrantOne()
		v++

		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 90)
		b.Variations[v].TranslateToQuadrantOne()
		v++
		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 180)
		b.Variations[v].TranslateToQuadrantOne()
		v++
		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 270)
		b.Variations[v].TranslateToQuadrantOne()
		v++
	}
	
	for omega := 90; omega < 360; omega += 180 {
		// should call for 90 and 270 (8 in total)
		o, _ = baseBlock.RotateAroundAxis(YAxis, omega)
		b.Variations[v] = o
		b.Variations[v].TranslateToQuadrantOne()
		v++

		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 90)
		b.Variations[v].TranslateToQuadrantOne()
		v++
		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 180)
		b.Variations[v].TranslateToQuadrantOne()
		v++
		b.Variations[v],_ = o.RotateAroundAxis(ZAxis, 270)
		b.Variations[v].TranslateToQuadrantOne()
		v++

	}


	return b

}
// a localization is a set of translations that will still keep the block within the block from 0,0,0 -> 3,3,3
func CreateBlockLocalizations(block *block) {

	for i:=0; i < len(block.Variations); i++) {
		block.Variations[i].CreateLocalizations()
	}

}



// This will serve any block
func (b *OrientedBlock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	blockJSON,_ := json.Marshal(*b)
	fmt.Printf("ServeHTTP serving block %s\n", blockJSON)
	fmt.Fprintf(w,"%s",blockJSON)

}

// calculates and fills in locations that will keep the oriented block, which is assumed to be in quadrant 0, 
// between 0,0,0 and 3,3,3
func (b *OrientedBlock) CreateLocalizations() {

	// block should be in quadrant 0 to start
	TODO, work here
	

}

// rotates block around axes and then brings it back into quadrant one
// only allow rotations of 90,180, or 270 degrees around an axis
func (b *OrientedBlock) RotateAroundAxis(axis Axis, degrees int) (*OrientedBlock, error) {
	switch degrees {
	case 90:
	case 180:
	case 270:
	default:
		return nil, errors.New("RotateAroundAxis can only rotate 90,180 and 270 degrees")
	}
	if (b == nil) {
		return nil, errors.New("RotateAroundAxis received nil for block")
	}

	newBlock := new(OrientedBlock)
	newBlock.BlockID = b.BlockID
	newBlock.Parts = make([]Coord, len(b.Parts))	// create a new array that will hold rotated points

	for i := 0; i < len(b.Parts); i++ {
		newBlock.Parts[i] = b.Parts[i].RotateAroundAxis(axis, degrees)
	}

	return newBlock, nil

	
}

// analyzes the block and translates it to Quadrant 1, putting it's lower left corner at the origin
func (b* OrientedBlock) TranslateToQuadrantOne() {

	// Analyze the block and figure out how to bring it to quadrant 1
	// Translate it to Quadrant 1

	b.TranslateToCoord(Coord{0,0,0})


}

// analyzes the block and translates it to specific coordinate, putting it's lower left corner at the coord
func (b* OrientedBlock) TranslateToCoord(c Coord) {
	// Analyze the block and figure out how to bring it to quadrant 1
	// Translate it to Quadrant 1
	var minx, miny, minz int = math.MaxInt32, math.MaxInt32, math.MaxInt32
	for i:=0; i < len(b.Parts); i++ {
		if b.Parts[i].X < minx {
			minx = b.Parts[i].X
		}
		if b.Parts[i].Y < miny {
			miny = b.Parts[i].Y
		}
		if b.Parts[i].Z < minz {
			minz = b.Parts[i].Z
		}
	}
	b.TranslateXYZ(c.X-minx, c.Y-miny, c.Z-minz)


}
// ignores blockID
// goes through every 
func (b* OrientedBlock) IsEqual(c* OrientedBlock) (bool) {
	if len(b.Parts) != len(c.Parts) {
		return false
	}
	for i := 0 ; i < len(b.Parts); i++ {
		
		if !(b.Parts[i].IsEqual(&(c.Parts[i]))) {
			return false
		}
	}
	return true
}


// this will translate every coordinate in the block by x,y,z
// it works in place
func (b *OrientedBlock) TranslateXYZ(x int, y int, z int) {

	// for each coordinate
	// TranslateCoord
	for i:=0; i < len(b.Parts); i++ {
		(&(b.Parts[i])).Translate(x, y, z)
	}

}
// convenience function
func (b *OrientedBlock) Translate(t Translation) {
	// for each coordinate
	// TranslateCoord
	for i:=0; i < len(b.Parts); i++ {
		(&(b.Parts[i])).Translate(t.X, t.Y, t.Z)
	}
}


// round a number
func Round(val float64) int {
    if val < 0 {
        return int(val-0.5)
    }
    return int(val+0.5)
}


// rotates a point around an axis
func (c Coord) RotateAroundAxis(axis Axis, degrees int) (newCoord Coord) {

	var radians float64 = float64(degrees)/360 * math.Pi * 2.0 // convert to radians
	var x,y,z float64
	var oldX, oldY, oldZ float64
	oldX = float64(c.X)
	oldY = float64(c.Y)
	oldZ = float64(c.Z)
	switch axis {
	case XAxis:
		y = oldY * math.Cos(radians) - oldZ * math.Sin(radians)
		z = oldY * math.Sin(radians) + oldZ * math.Cos(radians)
		x = oldX

	case YAxis:
		z = oldZ * math.Cos(radians) - oldX *math.Sin(radians)
		x = oldZ * math.Sin(radians) + oldX *math.Cos(radians)
		y = oldY

	case ZAxis:
		x = oldX * math.Cos(radians) - oldY * math.Sin(radians)
		y = oldX * math.Sin(radians) + oldY* math.Cos(radians)
		z = oldZ
	default:
		fmt.Println("Bad call to (Coord) RotateAroundAxis")
	}
	newCoord.X = Round(x)
	newCoord.Y = Round(y)
	newCoord.Z = Round(z)
	// fmt.Printf("Old Coord = %v, New Coord = %v\n", c, newCoord)
	return (newCoord)
}

// translates a coordinate 
func (c *Coord) Translate(x int, y int, z int) {
	c.X += x
	c.Y += y
	c.Z += z
}

func (c *Coord) IsEqual(d *Coord) (bool) {
	if c == nil || d == nil {
		return false
	}
	if (c.X != d.X) || (c.Y != d.Y) || (c.Z != d.Z) {
		return false
	}
	return true
}

// This is the primary solver
func Solve() (error) {
// outline
	// for each piece
	//   for each orientation
	//     for each starting postion of (-x,x, -y,y, -z, z)
	//        for each translation that still results in piece going into cube boundary
	//           translate block into final cube boundary in quadrant 1
	//           if block can't get into quadrant one, then continue


	var space *Space = NewSpace(SpaceDimension)
	for i:=0; i < len(blocks); i++ {
		for orient :=0; orient < len(blocks[i].Variations); orient++ {
			block := blocks[i].Variations[orient]
			
			for startingPos := 0; startingPos < len(startingPositions); startingPos++ {

				// Translate lower left corner of Block to that Location
				block.TranslateToCoord(startingPositions[startingPos])
				// try in this position
				// TODO put it in space

				success, e := space.InsertBlock(block)
				// this should ALWAYS work
				if (!success || e!= nil) {
					return errors.New("Solve() Internal error when trying to put oriented block in initial position")
				}

				// TODO now try to translate into cube position, if I fail, then this path is not a solution


				// need to apply local movements, which is the list of movements possible
				for perturb:=0; perturb < len(block.localMovements); perturb++ {
					block.Translate(block.localMovements[perturb])

					// try in this position

					success, e = space.InsertBlock(block)
					// this should ALWAYS work
					if (!success || e!= nil) {
						return errors.New("Solve() Internal error when trying to put oriented, perturbed block in initial position")
					}


					// TODO now try to translate it into cube position. If I fail then this path is not a solution


					block.TranslateToCoord(startingPositions[startingPos]) // reset it into starting position

				}

			}


		}
	}

	return nil

}

func main() {
	var block4JSON []byte
	var err error
	block4JSON, err = json.Marshal(block4)
	if (err == nil) {

		fmt.Printf("JSON for block4: %s\n",block4JSON)
	} else {
		fmt.Printf("error was %v",err)
	}


	// outline
	// create all orientations of blocks
	// for each piece
	//   for each orientation
	//     for each starting postion of (-x,x, -y,y, -z, z)
	//        for each translation that still results in piece going into cube boundary
	//           translate block into final cube boundary in quadrant 1
	//           if block can't get into quadrant one, then continue

	blocks[0] = CreateBlockOrientations(&block0)
	blocks[1] = CreateBlockOrientations(&block1)
	blocks[2] = CreateBlockOrientations(&block2)
	blocks[3] = CreateBlockOrientations(&block3)
	blocks[4] = CreateBlockOrientations(&block4)
	blocks[5] = CreateBlockOrientations(&block5)

	Solve()


	http.Handle("/block/0", &block0)
	http.Handle("/block/1", &block1)
	http.Handle("/block/2", &block2)
	http.Handle("/block/3", &block3)
	http.Handle("/block/4", blocks[4].Variations[19])
	http.Handle("/block/5", &block5)


	fs := http.FileServer(http.Dir("/Users/aje/js"))
	http.Handle("/",fs)
	fmt.Printf("About to listen on port 9000\n")
	fmt.Printf("%v",http.ListenAndServe(":9000",nil))
	fmt.Printf("Done listening\n")

}
