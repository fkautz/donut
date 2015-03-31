package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
	//	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func setupNodeDiskMap(c *C) map[string][]string {
	var disks []string
	for i := 0; i < 16; i++ {
		root, err := ioutil.TempDir(os.TempDir(), "donut-")
		c.Assert(err, IsNil)
		disks = append(disks, root)
	}
	nodeDiskMap := make(map[string][]string)
	nodeDiskMap["localhost"] = disks
	return nodeDiskMap
}

func removeDisks(c *C, disks []string) {
	for _, disk := range disks {
		err := os.RemoveAll(disk)
		c.Assert(err, IsNil)
	}
}

func (s *MySuite) TestEmptyBucket(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	// check buckets are empty
	buckets, err := donut.ListBuckets()
	c.Assert(err, IsNil)
	c.Assert(len(buckets), Equals, 0)
}

func (s *MySuite) TestBucketWithoutNameFails(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	// fail to create new bucket without a name
	err = donut.PutBucket("")
	c.Assert(err, Not(IsNil))

	err = donut.PutBucket(" ")
	c.Assert(err, Not(IsNil))
}

func (s *MySuite) TestCreateBucketAndList(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	// make bucket
	err = donut.PutBucket("foo")
	c.Assert(err, IsNil)

	// check bucket exists
	buckets, err := donut.ListBuckets()
	c.Assert(err, IsNil)
	c.Assert(len(buckets), Equals, 1)
	c.Assert(buckets[0].Name, Equals, "foo")
}

func (s *MySuite) TestCreateBucketWithSameNameFails(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	// make bucket
	err = donut.PutBucket("foo")
	c.Assert(err, IsNil)

	// make bucket fail
	err = donut.PutBucket("foo")
	c.Assert(err, Not(IsNil))
}

func getBucketNames(buckets []*Bucket) []string {
	var bucketNames []string
	for _, bucket := range buckets {
		bucketNames = append(bucketNames, bucket.Name)
	}
	return bucketNames
}

func (s *MySuite) TestCreateMultipleBucketsAndList(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	err = donut.PutBucket("foo")
	c.Assert(err, IsNil)

	err = donut.PutBucket("bar")
	c.Assert(err, IsNil)

	buckets, err := donut.ListBuckets()
	c.Assert(err, IsNil)

	bucketNamesReceived := getBucketNames(buckets)
	bucketNamesProvided := []string{"bar", "foo"}
	c.Assert(bucketNamesReceived, DeepEquals, bucketNamesProvided)

	err = donut.PutBucket("foobar")
	c.Assert(err, IsNil)
	bucketNamesProvided = append(bucketNamesProvided, "foobar")

	buckets, err = donut.ListBuckets()
	c.Assert(err, IsNil)
	bucketNamesReceived = getBucketNames(buckets)
	c.Assert(bucketNamesReceived, DeepEquals, bucketNamesProvided)
}

func (s *MySuite) TestNewObjectFailsWithoutBucket(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	_, _, err = donut.Get("foo", "obj")
	c.Assert(err, Not(IsNil))
}

func (s *MySuite) TestNewObjectFailsWithEmptyName(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	_, _, err = donut.Get("foo", "")
	c.Assert(err, Not(IsNil))

	_, _, err = donut.Get("foo", " ")
	c.Assert(err, Not(IsNil))
}

func (s *MySuite) TestNewObjectCanBeWritten(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	err = donut.PutBucket("foo")
	c.Assert(err, IsNil)

	data := "Hello World"
	err = donut.Put("foo", "obj", int64(len(data)), bytes.NewBuffer([]byte(data)))
	c.Assert(err, IsNil)

	reader, size, err := donut.Get("foo", "obj")
	c.Assert(err, IsNil)

	var actualData bytes.Buffer
	_, err = io.Copy(&actualData, reader)
	c.Assert(err, IsNil)
	c.Assert(actualData.Bytes(), DeepEquals, []byte(data))
	c.Assert(int64(actualData.Len()), Equals, size)
}

func getObjectNames(objects []*Item) []string {
	var objectNames []string
	for _, object := range objects {
		objectNames = append(objectNames, object.Key)
	}
	return objectNames
}

func (s *MySuite) TestMultipleNewObjects(c *C) {
	nodeDiskMap := setupNodeDiskMap(c)
	donut, err := GetNewClient("testemptydonut", nodeDiskMap)
	defer removeDisks(c, nodeDiskMap["localhost"])
	c.Assert(err, IsNil)

	c.Assert(donut.PutBucket("foo"), IsNil)
	err = donut.Put("foo", "obj1", int64(len([]byte("one"))), bytes.NewBuffer([]byte("one")))
	c.Assert(err, IsNil)

	err = donut.Put("foo", "obj2", int64(len([]byte("two"))), bytes.NewBuffer([]byte("two")))
	c.Assert(err, IsNil)

	reader, _, err := donut.Get("foo", "obj1")
	c.Assert(err, IsNil)

	var readerBuffer1 bytes.Buffer
	_, err = io.Copy(&readerBuffer1, reader)
	c.Assert(err, IsNil)
	c.Assert(readerBuffer1.Bytes(), DeepEquals, []byte("one"))

	reader, _, err = donut.Get("foo", "obj2")
	c.Assert(err, IsNil)

	var readerBuffer2 bytes.Buffer
	_, err = io.Copy(&readerBuffer2, reader)
	c.Assert(err, IsNil)
	c.Assert(readerBuffer2.Bytes(), DeepEquals, []byte("two"))

	// test list objects
	listObjects, _, err := donut.ListObjects("foo", "", "", "", Maxkeys)
	c.Assert(err, IsNil)

	receivedObjectNames := getObjectNames(listObjects)
	c.Assert(receivedObjectNames, DeepEquals, []string{"obj1", "obj2"})
}
