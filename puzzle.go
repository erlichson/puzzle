package main

import ("fmt" 
		"encoding/json"
		"net/http")

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




// Each block is described with a base orientation

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

// This will serve any block
func (b *OrientedBlock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	blockJSON,_ := json.Marshal(*b)
	fmt.Fprintf(w,"%s",blockJSON)

}

// rotates block around axes and then brings it back into quadrant one
// only allow rotations of 90,180, or 270 degrees around an axis
func (b *OrientedBlock) RotateAroundAxis(axis Axis, degrees int) (newBlock *OrientedBlock, error Error) {


	
}

// rotates a point around an axis
func (c Coord) RotateAroundAxis(axis Axis, degrees int) (newCoord Coord) {

	var radians float64 = degrees/360 * math.Pi * 2.0 // convert to radians
	var x,y,z float64
	var oldX, oldY, oldZ
	oldX = float64(c.X)
	oldY = float64(c.Y)
	oldZ = float64(c.Z)
	switch axis {
	case XAxis:
		y = oldY * math.Cos(q) - oldZ * math.Sin(q)
		z = oldY * math.Sin(q) + oldZ * math.Cos(q)
		x = oldX

	case YAxis:
		z = oldZ * math.Cos(q) - oldX *math.Sin(q)
		x = oldZ * math.Sin(q) + oldX *math.Cos(q)

	case ZAxis:
		x = oldX * math.Cos(q) - oldY * math.Sin(q)
		y = oldX * math.Sin(q) + oldY* math.Cos(q)
		z = oldZ
	default:
		fmt.Println("Bad call to (Coord) RotateAroundAxis")
	}


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

	// calculate some block rotations




	http.Handle("/block/0", &block0)
	http.Handle("/block/1", &block1)
	http.Handle("/block/2", &block2)
	http.Handle("/block/3", &block3)
	http.Handle("/block/4", &block4)	
	http.Handle("/block/5", &block5)


	fs := http.FileServer(http.Dir("/Users/aje/js"))
	http.Handle("/",fs)
	fmt.Printf("About to listen on port 9000\n")
	fmt.Printf("%v",http.ListenAndServe(":9000",nil))
	fmt.Printf("Done listening\n")

}
