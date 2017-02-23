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

type OrientedBlock struct {
	BlockID int	 `json:"blockID"`	// holds the ID of the block
	Parts []Coord `json:"parts"`	// holds a set of base coordinates of the pieces
}

type Block struct {
	BlockID int
	Variations [24]OrientedBlock
}

type Axis int
const (
	XAxis Axis = iota
	YAxis
	ZAxis)


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

// Given an oriented Block, return a Block with all 24 variations
func CreateBlockOrientations(baseBlock *OrientedBlock) (*Block) {
	b := new(Block)
	int v = 1
	b.BlockID = baseBlock.BlockID
	b.Variations[0] = baseBlock
	b.Variations[v++], _ = baseBlock.RotateAroundAxis(ZAxis, 90)
	b.Variations[v++], _ = baseBlock.RotateAroundAxis(Zaxis, 180)
	b.Variations[v++], _ = baseBlock.RotateAroundAxis(Zaxis, 270)

	b.Variations[v++], _ = baseBlock.RotateAroundAxis(Xaxis, 90)


	return b

}

// This will serve any block
func (b *OrientedBlock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	blockJSON,_ := json.Marshal(*b)
	fmt.Printf("ServeHTTP serving block %s\n", blockJSON)
	fmt.Fprintf(w,"%s",blockJSON)

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

// analyzes the block and translates it to Quadrant 1, putting it's lower left corner at the orgin
func (b* OrientedBlock) TranslateToQuadrantOne() {

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
	b.Translate(-minx, -miny, -minz)


}

// this will translate every coordinate in the block by x,y,z
// it works in place
func (b *OrientedBlock) Translate(x int, y int, z int) {

	// for each coordinate
	// TranslateCoord
	for i:=0; i < len(b.Parts); i++ {
		(&(b.Parts[i])).Translate(x, y, z)
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
	fmt.Printf("Old Coord = %v, New Coord = %v\n", c, newCoord)
	return (newCoord)
}

// translates a coordinate 
func (c *Coord) Translate(x int, y int, z int) {
	c.X += x
	c.Y += y
	c.Z += z
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


	http.Handle("/block/0", &block0)
	http.Handle("/block/1", &block1)
	http.Handle("/block/2", &block2)
	http.Handle("/block/3", &block3)
    rotatedBlock4,_ := block4.RotateAroundAxis(ZAxis,90)
    rotatedBlock4.TranslateToQuadrantOne()
	//http.Handle("/block/4", &block4)	
	http.Handle("/block/4", rotatedBlock4)
	http.Handle("/block/5", &block5)


	fs := http.FileServer(http.Dir("/Users/aje/js"))
	http.Handle("/",fs)
	fmt.Printf("About to listen on port 9000\n")
	fmt.Printf("%v",http.ListenAndServe(":9000",nil))
	fmt.Printf("Done listening\n")

}
