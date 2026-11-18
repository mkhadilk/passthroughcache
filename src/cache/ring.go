package cache

import (
	"container/ring"
	"errors"
	"fmt"
)

type requestType int

const (
	getRequest requestType = iota
	refreshRequest
)

type Ideable interface {
	GetId() string
}

type reQuest[T any] struct {
	atype requestType
	id    string
	resp  chan T
}

func (r *reQuest[T]) readResponse() (res T, err error) {
	rs, ok := <-r.resp
	if ok {
		return rs, nil
	} else {
		return rs, fmt.Errorf("Can not get response for entity id: %s", r.id)
	}
}

type ARing[T any] struct {
	name           string
	requestChannel chan reQuest[T]
	buffersize     int
	fetchHandler   func(id string) (Ideable, error)
}

func NewRing[Etype any](name string, buffersize int) ARing[Etype] {

	var ar ARing[Etype]
	ar.name = name
	ar.buffersize = buffersize
	return ar
}

func (a *ARing[T]) Stop() {
	if a.requestChannel != nil {
		close(a.requestChannel)
		a.requestChannel = nil
	}
}
func (a *ARing[T]) Start() {
	if a.requestChannel == nil {
		aringchan := make(chan reQuest[T], 10)
		a.requestChannel = aringchan
		go a.process()
	}
}
func (a *ARing[T]) process() {
	aringbuffer := ring.New(a.buffersize)
	for {
		if rr, ok := <-a.requestChannel; ok {

			vl, err := a.locateInRing(aringbuffer, rr.id)

			switch rr.atype {
			case getRequest:
				if err != nil {
					t, ferr := a.fetchHandler(rr.id)
					if ferr == nil {
						aringbuffer.Value = t.(T)
						vl = t.(T)
					} else {
						aringbuffer.Value = nil
					}
				}
				rr.resp <- vl
				aringbuffer = aringbuffer.Next()
			case refreshRequest:
				t, ferr := a.fetchHandler(rr.id)
				if ferr == nil {
					aringbuffer.Value = t
				} else {
					aringbuffer.Value = nil
				}
				gotit := aringbuffer.Value.(T)
				rr.resp <- gotit
			}
		} else {
			break
		}
	}
}

func (a *ARing[T]) GetName() string {
	return a.name
}

func (a *ARing[T]) locateInRing(aring *ring.Ring, id string) (T, error) {
	var t T
	for i := 0; i < aring.Len(); i++ {
		if aring.Value != nil {
			idb, ok := aring.Value.(Ideable)
			if ok {
				if idb.GetId() == id {
					gotit := aring.Value.(T)
					return gotit, nil
				}
			}

		}
		aring = aring.Next()
	}
	return t, errors.New("Not found")
}

func (a *ARing[T]) Get(id string) (res T, err error) {
	if a.requestChannel != nil {
		var newr reQuest[T]
		newr.atype = getRequest
		newr.id = id
		newr.resp = make(chan T, 1)
		a.requestChannel <- newr
		res, err = newr.readResponse()
	} else {
		err = errors.New("Not started")
	}
	return
}
func (a *ARing[T]) Refresh(id string) {
	if a.requestChannel != nil {
		var newr reQuest[T]
		newr.atype = refreshRequest
		newr.id = id
		newr.resp = make(chan T, 1)
		a.requestChannel <- newr
		newr.readResponse()
	}
	return
}

func (a *ARing[T]) SetFetchHandler(getFunc func(id string) (Ideable, error)) {
	a.fetchHandler = getFunc
}
