package page

import "github.com/rivo/tview"

type Layout struct {
	nav    map[string]*Position
	width  int
	height int
}

type Position struct {
	Id     string
	Pos    [2]int
	Height int
	Width  int
	Left   string
	Right  string
	Up     string
	Down   string
}

type RowCol struct {
	row []int
	col []int
}

func (rc *RowCol) setRow(ele int) bool {
	for _, a := range rc.row {
		if a == ele {
			return false
		}
	}

	rc.row = append(rc.row, ele)
	return true
}

func (rc *RowCol) setCol(ele int) bool {
	for _, a := range rc.col {
		if a == ele {
			return false
		}
	}

	rc.col = append(rc.col, ele)
	return true
}

func MakeLayout(nav [][]string) *Layout {
	positions := map[string]*Position{}
	rowCol := map[string]*RowCol{}

	for i, row := range nav {
		for j, name := range row {
			if name == "" {
				continue
			}
			pos := positions[name]
			rc := rowCol[name]
			if pos == nil {
				pos = &Position{Id: name, Pos: [2]int{i, j}, Height: 1, Width: 1}
				rc = &RowCol{row: []int{i}, col: []int{j}}
				positions[name] = pos
				rowCol[name] = rc
			} else {
				if rc.setRow(i) {
					pos.Height++
				}
				if rc.setCol(j) {
					pos.Width++
				}
			}

			if j > 0 && pos.Left == "" {
				for k := j - 1; k >= 0; k-- {
					left := nav[i][k]
					if left != name && left != "" {
						pos.Left = left
						break
					}
				}
			}

			if j < len(row)-1 && pos.Right == "" {
				for k := j + 1; k < len(row); k++ {
					right := nav[i][k]
					if right != name && right != "" {
						pos.Right = right
						break
					}
				}
			}

			if i > 0 && pos.Up == "" {
				for k := i - 1; k >= 0; k-- {
					up := nav[k][j]
					if up != name && up != "" {
						pos.Up = up
						break
					}
				}
			}

			if i < len(nav)-1 && pos.Down == "" {
				for k := i + 1; k < len(nav); k++ {
					down := nav[k][j]
					if down != name && down != "" {
						pos.Down = down
						break
					}
				}
			}
		}
	}

	return &Layout{
		nav:    positions,
		width:  len(nav[0]),
		height: len(nav),
	}
}

func (nav *Layout) Setup(grid *tview.Grid) {
	rows := make([]int, nav.height)
	cols := make([]int, nav.width)

	grid.SetRows(rows...).SetColumns(cols...)
}

func (nav *Layout) GetIds() []string {
	keys := make([]string, 0, len(nav.nav))
	for name := range nav.nav {
		keys = append(keys, name)
	}

	return keys
}

func (nav *Layout) GetPos(position string) *Position {
	return nav.nav[position]
}

func (nav *Layout) Left(position string) string {
	return nav.nav[position].Left
}

func (nav *Layout) Right(position string) string {
	return nav.nav[position].Right
}

func (nav *Layout) Up(position string) string {
	return nav.nav[position].Up
}

func (nav *Layout) Down(position string) string {
	return nav.nav[position].Down
}
