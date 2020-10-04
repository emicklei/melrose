module github.com/emicklei/melrose/cmd/termrose

go 1.14

require (
	github.com/gdamore/tcell v1.4.0
	github.com/rivo/tview v0.0.0-20200915114512-42866ecf6ca6
	github.com/emicklei/melrose v0.30.0
)

replace github.com/emicklei/melrose v0.30.0 => ../..
