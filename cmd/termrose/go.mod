module github.com/emicklei/melrose/cmd/termrose

go 1.14

require (
	github.com/emicklei/melrose v0.30.0
	github.com/gdamore/tcell v1.4.0
	github.com/gdamore/tcell/v2 v2.0.0-dev // indirect
	github.com/rivo/tview v0.0.0-20200915114512-42866ecf6ca6
)

replace github.com/emicklei/melrose v0.30.0 => ../..
