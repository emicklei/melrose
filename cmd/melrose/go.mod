module github.com/emicklei/melrose/cmd/melrose

go 1.14

require (
    github.com/emicklei/melrose v0.3.1
    github.com/peterh/liner v1.2.0
)

replace github.com/emicklei/melrose => ../../
replace github.com/peterh/liner => ../../../liner