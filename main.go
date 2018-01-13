package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

// Bash: clear && go build -o binary.exe && binary.exe maze1
// PowerShell (integrated VSCode terminal): clear; go build -o binary.exe; ./binary.exe maze1

func main() {
	if len(os.Args) < 2 {
		log.Println("No path to maze image")
		log.Println("Exiting...")
		os.Exit(1)
	}

	// Prepare solution folders if not exists

	if _, err := os.Stat("./solutions"); os.IsNotExist(err) {
		fError := os.Mkdir("solutions", os.ModePerm)
		if fError != nil {
			log.Fatal("Setup error:", fError)
		}
	}

	// Step 1: Load the image file
	path := os.Args[1]
	imagePath := fmt.Sprintf("./mazes/%v.png", path)
	fReader, err := os.Open(imagePath)
	defer fReader.Close()
	if err != nil {
		log.Fatalf("Error opening file. Error: %v", err)
	}

	// Step 2: Convert image into a mesh of Nodes
	img, _, err := image.Decode(fReader)
	if err != nil {
		log.Fatalf("Error decoding file. Error: %v", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	log.Printf("Image is %v x %v", width, height)

	var nodes []Node

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var node Node
			tempR, tempG, tempB, _ := img.At(x, y).RGBA()             //uint32 values (Skip alpha value)
			r, g, b := int(tempR/257), int(tempG/257), int(tempB/257) // Somehow convert uint32 to 0-255 RGB color scale (Stackoverflow is great)
			if r == 255 && g == 255 && b == 255 {
				node = Node{isWall: false}
			} else {
				node = Node{isWall: true}
			}

			nodes = append(nodes, node)
		}
	}

	for i := range nodes {
		if i > 0 {
			nodes[i].left = &nodes[i-1]
		}

		if i+1 < len(nodes)-1 {
			nodes[i].right = &nodes[i+1]
		}

		if i-width > 0 {
			nodes[i].up = &nodes[i-width]
		}

		if i+width < len(nodes)-1 {
			nodes[i].down = &nodes[i+width]
		}

		nodes[i].checked = false
	}

	var startPos *Node
	var endPos *Node

	for i := 0; i < width; i++ {
		if nodes[i].isWall == false {
			startPos = &nodes[i]
			log.Println("Got starting position at:", i)
		}
	}

	for i := len(nodes) - 1 - width; i < len(nodes)-1; i++ {
		if nodes[i].isWall == false {
			endPos = &nodes[i]
			log.Println("Got end pos at:", i)
		}
	}

	// log.Println(startPos, endPos, len(nodes)-1)

	startPos.checked = true
	stack := NewStack()
	stack.Push(startPos.down)

	parent := startPos

	for {
		n := stack.Pop()
		if n == nil {
			// End of the stack, either errored or finished
			log.Println("Ran out of items")
			return
		}

		n.checked = true  // Make sure current item wont be processed again
		n.parent = parent // Make sure item finds it way back
		parent = n        // Set new parent

		if n == endPos {
			log.Println("Got to the end")
			log.Println("n:", n)
			str, err := OutputSolution(path, height, width)
			if err != nil {
				log.Fatal("Error1:", err)
			}
			log.Println("OutputSolution returned:", str)
			return
		}

		if n.up != nil {
			if n.up.checked == false {
				stack.Push(n.up)
			}
		}

		if n.down != nil {
			if n.down.checked == false {
				stack.Push(n.down)
			}
		}

		if n.left != nil {
			if n.left.checked == false {
				stack.Push(n.left)
			}
		}

		if n.right != nil {
			if n.right.checked == false {
				stack.Push(n.right)
			}
		}
	}
}