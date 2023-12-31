package kademlia

import (
	"testing"
)

func TestBucket(t *testing.T) {
	bucket := newBucket()
	mainContact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	bucket.AddContact(mainContact)
	bucket.AddContact(NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:8002"))
	bucket.AddContact(mainContact)
	if bucket.Len() == 3 {
		t.Errorf("bucket was able to add a duplicate!")
	}
	//try to overload a bucket
	for i := 0; i < bucketSize+5; i++ {
		bucket.AddContact(NewContact(NewRandomKademliaID(), "localhost:8001"))
	}
	if bucket.Len() > bucketSize {
		t.Errorf("bucket was bigger (%d) than max allowed (%d)", bucket.Len(), bucketSize)
	}

}
