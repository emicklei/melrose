package main

var refmgt = NewReferenceManager()

type RefMgr struct {
	objects map[int]interface{}
}

func NewReferenceManager() *RefMgr {
	objs := map[int]interface{}{}
	return &RefMgr{objs}
}
