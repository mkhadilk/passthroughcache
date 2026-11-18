package cache

import (
	"crypto/rand"
	"errors"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

type localentity struct {
	Ideable
	ID  string
	Val string
}

func (l localentity) GetId() string {
	return l.ID
}

func TestCreateRing(t *testing.T) {
	r := NewRing[localentity]("Some Name", 2)

	log.Printf("Ring Name %s", r.GetName())

	defer r.Stop()

	if r.GetName() != "Some Name" {
		log.Printf("Wrong Ring Name set; expecting %s; found %s", "Some Name", r.GetName())
		t.FailNow()
	}
}
func TestCreateRingSize(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	defer r.Stop()

	if r.buffersize != 12 {
		log.Printf("Wrong Ring Size set; expecting %d; found %d", 12, r.buffersize)
		t.FailNow()
	}
}

func TestBlankRingFetchNoStart(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	r.SetFetchHandler(func(id string) (val Ideable, err error) {
		if len(id) > 15 {
			val = nil
			err = errors.New("Not found")
		} else {
			var vall localentity

			vall.ID = "blah"
			val = vall
			err = nil
		}
		return
	})

	_, err := r.Get("Bah")
	if err == nil {
		log.Println("Missing not started error")
		t.FailNow()
	}

}
func TestBlankRingFetch(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	r.SetFetchHandler(func(id string) (val Ideable, err error) {
		if len(id) > 15 {
			val = nil
			err = errors.New("Not found")
		} else {
			var vall localentity

			vall.ID = id

			val = vall
			err = nil
		}
		return
	})

	r.Start()

	defer r.Stop()

	_, err := r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

}

func TestBlankRingFetchRing(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	r.SetFetchHandler(func(id string) (val Ideable, err error) {
		if len(id) > 15 {
			val = nil
			err = errors.New("Not found")
		} else {
			var vall localentity

			vall.ID = id

			val = vall
			err = nil
		}
		return
	})

	r.Start()

	defer r.Stop()

	_, err := r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

	_, err = r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

}

func TestBlankRingFetchRingLong(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	r.SetFetchHandler(func(id string) (val Ideable, err error) {
		if len(id) > 15 {
			val = nil
			err = errors.New("Not found")
		} else {
			var vall localentity

			vall.ID = id

			val = vall
			err = nil
		}
		return
	})

	r.Start()

	defer r.Stop()

	_, err := r.Get("Bahqqqwwqeqwewqewqewqewqewqe")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

	_, err = r.Get("Bahqqqwwqeqwewqewqewqewqewqe")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

}

func TestBlankRingRefreshRing(t *testing.T) {
	r := NewRing[localentity]("Some Name", 12)

	log.Printf("Ring Size %d", r.buffersize)

	r.SetFetchHandler(func(id string) (val Ideable, err error) {
		if len(id) <= 5 {
			var vall localentity

			vall.ID = id
			vall.Val = rand.Text()

			val = vall
			err = nil
		} else if len(id) > 35 {
			val = nil
			err = errors.New("Not found")
		} else {
			var vall localentity

			vall.ID = id
			vall.Val = rand.Text()

			val = vall
			err = nil
		}
		return
	})

	r.Start()

	defer r.Stop()

	vv, err := r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

	log.Printf("Got Value %+v", vv)

	vl, err := r.Get("Bahhhhhhhhhhhhhhh")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}

	log.Printf("Got Value %+v", vl)

	vl, err = r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}
	log.Printf("Got Value before refresh %+v", vl)

	r.Refresh("Bah")
	vl, err = r.Get("Bah")
	if err != nil {
		log.Printf("Get failed %+v", err.Error())
		t.FailNow()
	}
	log.Printf("Got Value after refresh %+v", vl)

}
