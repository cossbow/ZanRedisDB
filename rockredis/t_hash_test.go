package rockredis

import (
	"bytes"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/absolute8511/ZanRedisDB/common"
	"github.com/stretchr/testify/assert"
)

func TestHashCodec(t *testing.T) {
	key := []byte("key")
	field := []byte("field")

	ek := hEncodeSizeKey(key)
	if k, err := hDecodeSizeKey(ek); err != nil {
		t.Fatal(err)
	} else if string(k) != "key" {
		t.Fatal(string(k))
	}

	ek = hEncodeHashKey([]byte("test"), key, field)
	if table, k, f, err := hDecodeHashKey(ek); err != nil {
		t.Fatal(err)
	} else if string(k) != "key" {
		t.Fatal(string(k))
	} else if string(f) != "field" {
		t.Fatal(string(f))
	} else if string(table) != "test" {
		t.Fatal(string(table))
	}
}

func TestDBHash(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	key := []byte("test:testdb_hash_a")

	if n, err := db.HSet(0, key, []byte("a"), []byte("hello world 1")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	if n, err := db.HSet(0, key, []byte("b"), []byte("hello world 2")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	v1, _ := db.HGet(key, []byte("a"))
	v2, _ := db.HGet(key, []byte("b"))
	if string(v1) != "hello world 1" {
		t.Error(v1)
	}
	if string(v2) != "hello world 2" {
		t.Error(v2)
	}
	ay, _ := db.HMget(key, []byte("a"), []byte("b"))

	if string(v1) != string(ay[0]) {
		t.Error(ay[0])
	}

	if string(v2) != string(ay[1]) {
		t.Error(ay[1])
	}

	len, err := db.HLen(key)
	if err != nil {
		t.Error(err)
	}
	if len != 2 {
		t.Errorf("length should be 2: %v", len)
	}
	_, ch, _ := db.HGetAll(key)
	results := make([]common.KVRecordRet, 0)
	for r := range ch {
		results = append(results, r)
	}
	if string(results[0].Rec.Key) != "a" {
		t.Error(results)
	}
	if string(results[0].Rec.Value) != "hello world 1" {
		t.Error(results)
	}
	if string(results[1].Rec.Key) != "b" {
		t.Error(results)
	}

	if string(results[1].Rec.Value) != "hello world 2" {
		t.Error(results)
	}

	_, ch, _ = db.HKeys(key)
	results = make([]common.KVRecordRet, 0)
	for r := range ch {
		results = append(results, r)
	}
	if string(results[0].Rec.Key) != "a" {
		t.Error(results)
	}
	if string(results[1].Rec.Key) != "b" {
		t.Error(results)
	}
	_, ch, _ = db.HValues(key)
	results = make([]common.KVRecordRet, 0)
	for r := range ch {
		results = append(results, r)
	}
	if string(results[0].Rec.Value) != "hello world 1" {
		t.Error(results)
	}
	if string(results[1].Rec.Value) != "hello world 2" {
		t.Error(results)
	}
}

func TestHashKeyExists(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()
	key := []byte("test:hkeyexists_test")
	v, err := db.HKeyExists(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != 0 {
		t.Fatal("invalid value ", v)
	}

	if _, err := db.HSet(0, key, []byte("hello"), []byte("world")); err != nil {
		t.Fatal(err.Error())
	}

	v, err = db.HKeyExists(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != 1 {
		t.Fatal("invalid value ", v)
	}
	if _, err := db.HSet(0, key, []byte("hello2"), []byte("world2")); err != nil {
		t.Fatal(err.Error())
	}
	db.HDel(key, []byte("hello"))
	v, err = db.HKeyExists(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != 1 {
		t.Fatal("invalid value ", v)
	}
	db.HClear(key)
	v, err = db.HKeyExists(key)
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != 0 {
		t.Fatal("invalid value ", v)
	}
}

func TestHashKeyIncrBy(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()
	key := []byte("test:hkey_incr_test")
	if _, err := db.HSet(0, key, []byte("hello"), []byte("0")); err != nil {
		t.Fatal(err.Error())
	}

	r, _ := db.HIncrBy(0, key, []byte("hello"), 3)
	if r != 3 {
		t.Error(r)
	}
	r, _ = db.HIncrBy(0, key, []byte("hello"), -6)
	if r != -3 {
		t.Error(r)
	}
}

func TestHashIndexLoad(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	assert.Equal(t, 0, len(db.indexMgr.tableIndexes))
	var hindex HsetIndex
	hindex.Table = []byte("test_index_table")
	hindex.Name = []byte("index1")
	hindex.IndexField = []byte("index_test_field")
	hindex.Unique = 0
	hindex.ValueType = StringV

	err := db.indexMgr.AddHsetIndex(db, &hindex)
	assert.Nil(t, err)

	var hindex2 HsetIndex
	hindex2.Table = []byte("test_index_table")
	hindex2.Name = []byte("index2")
	hindex2.IndexField = []byte("index_test_field2")
	hindex2.Unique = 0
	hindex2.ValueType = Int64V

	err = db.indexMgr.AddHsetIndex(db, &hindex2)
	assert.Nil(t, err)

	var hindex3 HsetIndex
	hindex3.Table = []byte("test_index_table3")
	hindex3.Name = []byte("index3")
	hindex3.IndexField = []byte("index_test_field3")
	hindex3.Unique = 0
	hindex3.ValueType = Int64V

	err = db.indexMgr.AddHsetIndex(db, &hindex3)
	assert.Nil(t, err)

	db.indexMgr.Close()
	err = db.indexMgr.LoadIndexes(db)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(db.indexMgr.tableIndexes))
	tindexes, err := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
	assert.Nil(t, err)
	assert.Equal(t, hindex.Name, tindexes.Name)

	tindexes, err = db.indexMgr.GetHsetIndex(string(hindex2.Table), string(hindex2.IndexField))
	assert.Nil(t, err)
	assert.Equal(t, hindex2.Name, tindexes.Name)

	tindexes, err = db.indexMgr.GetHsetIndex(string(hindex3.Table), string(hindex3.IndexField))
	assert.Nil(t, err)
	assert.Equal(t, hindex3.Name, tindexes.Name)
}

func TestHashIndexBuildAndClean(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	assert.Equal(t, 0, len(db.indexMgr.tableIndexes))
	var hindex HsetIndex
	hindex.Table = []byte("test_index_table")
	hindex.Name = []byte("index1")
	hindex.IndexField = []byte("index_test_field")
	hindex.Unique = 0
	hindex.ValueType = StringV

	key := []byte(string(hindex.Table) + ":testdb_hash_a")
	db.HSet(0, key, hindex.IndexField, []byte("1"))
	key = []byte(string(hindex.Table) + ":testdb_hash_b")
	db.HSet(0, key, hindex.IndexField, []byte("2"))
	key = []byte(string(hindex.Table) + ":testdb_hash_c")
	db.HSet(0, key, hindex.IndexField, []byte("2"))

	err := db.indexMgr.AddHsetIndex(db, &hindex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex.Table), string(hindex.IndexField), BuildingIndex)
	assert.Nil(t, err)
	// wait until building done
	buildStart := time.Now()
	for {
		time.Sleep(time.Millisecond * 10)
		hindex, err := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
		assert.Nil(t, err)
		if hindex.State == BuildDoneIndex {
			break
		} else if time.Since(buildStart) > time.Second*10 {
			t.Errorf("building index timeout")
			break
		}
	}
	condAll := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	cnt, _, err := db.HsetIndexSearch(hindex.Table, hindex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(cnt))

	var hindex2 HsetIndex
	hindex2.Table = []byte("test_index_table")
	hindex2.Name = []byte("index2")
	hindex2.IndexField = []byte("index_test_field2")
	hindex2.Unique = 0
	hindex2.ValueType = Int64V

	err = db.indexMgr.AddHsetIndex(db, &hindex2)
	assert.Nil(t, err)

	var hindex3 HsetIndex
	hindex3.Table = []byte("test_index_table3")
	hindex3.Name = []byte("index3")
	hindex3.IndexField = []byte("index_test_field3")
	hindex3.Unique = 0
	hindex3.ValueType = Int64V

	err = db.indexMgr.AddHsetIndex(db, &hindex3)
	assert.Nil(t, err)

	key = []byte(string(hindex2.Table) + ":testdb_hash_a")
	db.HSet(0, key, hindex2.IndexField, []byte("1"))
	key = []byte(string(hindex2.Table) + ":testdb_hash_b")
	db.HSet(0, key, hindex2.IndexField, []byte("2"))
	key = []byte(string(hindex2.Table) + ":testdb_hash_c")
	db.HSet(0, key, hindex2.IndexField, []byte("2"))

	key = []byte(string(hindex3.Table) + ":testdb_hash_a")
	db.HSet(0, key, hindex3.IndexField, []byte("1"))
	key = []byte(string(hindex3.Table) + ":testdb_hash_b")
	db.HSet(0, key, hindex3.IndexField, []byte("2"))
	key = []byte(string(hindex3.Table) + ":testdb_hash_c")
	db.HSet(0, key, hindex3.IndexField, []byte("2"))
	key = []byte(string(hindex3.Table) + ":testdb_hash_d")
	db.HSet(0, key, hindex3.IndexField, []byte("3"))

	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex2.Table), string(hindex2.IndexField), BuildingIndex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex3.Table), string(hindex3.IndexField), BuildingIndex)
	assert.Nil(t, err)
	// wait until building done
	buildStart = time.Now()
	for {
		time.Sleep(time.Millisecond * 10)
		hindex2, err := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
		assert.Nil(t, err)
		hindex3, err := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
		assert.Nil(t, err)
		if hindex2.State == BuildDoneIndex && hindex3.State == BuildDoneIndex {
			break
		} else if time.Since(buildStart) > time.Second*10 {
			t.Errorf("building index timeout")
			break
		}
	}

	cnt, _, err = db.HsetIndexSearch(hindex2.Table, hindex2.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(cnt))

	cnt, _, err = db.HsetIndexSearch(hindex3.Table, hindex3.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, 4, int(cnt))

	// clean index
	deletedIndex1, err := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
	assert.Nil(t, err)
	deletedIndex2, err := db.indexMgr.GetHsetIndex(string(hindex2.Table), string(hindex2.IndexField))
	assert.Nil(t, err)
	deletedIndex3, err := db.indexMgr.GetHsetIndex(string(hindex3.Table), string(hindex3.IndexField))
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex.Table), string(hindex.IndexField), DeletedIndex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex2.Table), string(hindex2.IndexField), DeletedIndex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex3.Table), string(hindex3.IndexField), DeletedIndex)
	assert.Nil(t, err)

	buildStart = time.Now()
	for {
		time.Sleep(time.Millisecond * 10)
		_, err1 := db.indexMgr.GetHsetIndex(string(hindex.Table), string(hindex.IndexField))
		_, err2 := db.indexMgr.GetHsetIndex(string(hindex2.Table), string(hindex2.IndexField))
		_, err3 := db.indexMgr.GetHsetIndex(string(hindex3.Table), string(hindex3.IndexField))
		if err1 == ErrIndexNotExist && err2 == ErrIndexNotExist && err3 == ErrIndexNotExist {
			break
		} else if time.Since(buildStart) > time.Second*10 {
			t.Errorf("clean index timeout")
			break
		}
	}

	cnt, _, err = deletedIndex1.SearchRec(db, condAll, false)
	assert.Equal(t, 0, int(cnt))

	cnt, _, err = deletedIndex2.SearchRec(db, condAll, false)
	assert.Equal(t, 0, int(cnt))

	cnt, _, err = deletedIndex3.SearchRec(db, condAll, false)
	assert.Equal(t, 0, int(cnt))
}

func TestHashIndexStringV(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	var hindex HsetIndex
	hindex.Table = []byte("test")
	hindex.Name = []byte("index1")
	hindex.IndexField = []byte("index_test_stringfield")
	hindex.Unique = 0
	hindex.ValueType = StringV

	err := db.indexMgr.AddHsetIndex(db, &hindex)
	assert.Nil(t, err)

	inputPKList := make([][]byte, 0, 3)
	inputPKList = append(inputPKList, []byte("test:key1"))
	inputPKList = append(inputPKList, []byte("test:key2"))
	inputPKList = append(inputPKList, []byte("test:key3"))
	inputFVList := make([][]byte, 0, 3)
	inputFVList = append(inputFVList, []byte("fv1"))
	inputFVList = append(inputFVList, []byte("fv2"))
	inputFVList = append(inputFVList, []byte("fv3"))
	db.wb.Clear()
	for i, pk := range inputPKList {
		err = db.hsetIndexAddRec(pk, hindex.IndexField, inputFVList[i], db.wb)
		assert.Nil(t, err)
	}
	db.eng.Write(db.defaultWriteOpts, db.wb)
	condAll := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}

	condEqual := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: true,
		EndKey:       inputFVList[0],
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}
	condLess := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       inputFVList[len(inputFVList)-1],
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	condLessEq := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       inputFVList[len(inputFVList)-1],
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}

	condGt := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	condGtEq := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: true,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	cnt, pkList, err := db.HsetIndexSearch(hindex.Table, hindex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condEqual, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, inputPKList[0], pkList[0].PKey)
	assert.Equal(t, inputFVList[0], pkList[0].IndexValue)

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLess, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
		assert.Equal(t, 1, bytes.Compare(condLess.EndKey, pkList[i].IndexValue))
	}
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLessEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
		comp := bytes.Compare(condLess.EndKey, pkList[i].IndexValue)
		assert.True(t, comp == 1 || comp == 0)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGt, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGtEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}

	db.wb.Clear()
	db.hsetIndexRemoveRec(inputPKList[0], hindex.IndexField, inputFVList[0], db.wb)
	db.eng.Write(db.defaultWriteOpts, db.wb)
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condEqual, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(cnt))

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLess, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-2, int(cnt))
	for i := 0; i < len(inputPKList)-2; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLessEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGt, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGtEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}
}

func TestHashIndexStringVPrefix(t *testing.T) {
}

func TestHashIndexStringVMultiEqual(t *testing.T) {
}

func TestHashIndexStringVUnique(t *testing.T) {
}

func TestHashIndexInt64V(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	var hindex HsetIndex
	hindex.Table = []byte("test")
	hindex.Name = []byte("index1")
	hindex.IndexField = []byte("index_test_int64field")
	hindex.Unique = 0
	hindex.ValueType = Int64V

	err := db.indexMgr.AddHsetIndex(db, &hindex)
	assert.Nil(t, err)

	inputPKList := make([][]byte, 0, 3)
	inputPKList = append(inputPKList, []byte("test:key1"))
	inputPKList = append(inputPKList, []byte("test:key2"))
	inputPKList = append(inputPKList, []byte("test:key3"))
	inputFVList := make([][]byte, 0, 3)
	inputFVList = append(inputFVList, []byte("1"))
	inputFVList = append(inputFVList, []byte("2"))
	inputFVList = append(inputFVList, []byte("11"))
	db.wb.Clear()
	for i, pk := range inputPKList {
		db.hsetIndexAddRec(pk, hindex.IndexField, inputFVList[i], db.wb)
	}
	db.eng.Write(db.defaultWriteOpts, db.wb)
	condAll := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	condEqual := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: true,
		EndKey:       inputFVList[0],
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}
	condLess := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       inputFVList[len(inputFVList)-1],
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	condLessEq := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       inputFVList[len(inputFVList)-1],
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}

	condGt := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	condGtEq := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: true,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}

	cnt, pkList, err := db.HsetIndexSearch(hindex.Table, hindex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}
	t.Log(pkList)

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condEqual, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, inputPKList[0], pkList[0].PKey)
	assert.Equal(t, inputFVList[0], pkList[0].IndexValue)

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLess, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLessEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGt, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGtEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))
	for i := 0; i < len(inputPKList); i++ {
		assert.Equal(t, inputPKList[i], pkList[i].PKey)
	}

	db.wb.Clear()
	db.hsetIndexRemoveRec(inputPKList[0], hindex.IndexField, inputFVList[0], db.wb)
	db.eng.Write(db.defaultWriteOpts, db.wb)
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condEqual, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(cnt))

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLess, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-2, int(cnt))
	for i := 0; i < len(inputPKList)-2; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}
	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condLessEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGt, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}

	cnt, pkList, err = db.HsetIndexSearch(hindex.Table, hindex.IndexField, condGtEq, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))
	for i := 0; i < len(inputPKList)-1; i++ {
		assert.Equal(t, inputPKList[i+1], pkList[i].PKey)
	}
}

func TestHashUpdateWithIndex(t *testing.T) {
	db := getTestDB(t)
	defer os.RemoveAll(db.cfg.DataDir)
	defer db.Close()

	var hindex HsetIndex
	hindex.Table = []byte("test")
	hindex.Name = []byte("index1")
	hindex.IndexField = []byte("index_test_int64field")
	hindex.Unique = 0
	hindex.ValueType = Int64V

	intIndex := hindex
	err := db.indexMgr.AddHsetIndex(db, &intIndex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex.Table), string(hindex.IndexField), ReadyIndex)
	assert.Nil(t, err)

	hindex.Table = []byte("test")
	hindex.Name = []byte("index2")
	hindex.IndexField = []byte("index_test_stringfield")
	hindex.Unique = 0
	hindex.ValueType = StringV

	stringIndex := hindex
	err = db.indexMgr.AddHsetIndex(db, &stringIndex)
	assert.Nil(t, err)
	err = db.indexMgr.UpdateHsetIndexState(db, string(hindex.Table), string(hindex.IndexField), ReadyIndex)
	assert.Nil(t, err)

	inputPKList := make([][]byte, 0, 3)
	inputPKList = append(inputPKList, []byte("test:testhindex_key1"))
	inputPKList = append(inputPKList, []byte("test:testhindex_key2"))
	inputPKList = append(inputPKList, []byte("test:testhindex_key3"))

	for i, key := range inputPKList {
		v := strconv.Itoa(i + 1)
		db.HSet(0, key, intIndex.IndexField, []byte(v))
	}

	condAll := &IndexCondition{
		StartKey:     nil,
		IncludeStart: false,
		EndKey:       nil,
		IncludeEnd:   false,
		Offset:       0,
		Limit:        -1,
	}
	cnt, _, err := db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))

	condEqual2 := &IndexCondition{
		StartKey:     []byte("2"),
		IncludeStart: true,
		EndKey:       []byte("2"),
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}

	cnt, pkList, err := db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condEqual2, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, 1, len(pkList))
	assert.Equal(t, inputPKList[1], pkList[0].PKey)
	assert.Equal(t, condEqual2.StartKey, pkList[0].IndexValue)
	// update field to new and check index update
	db.HSet(0, inputPKList[1], intIndex.IndexField, []byte("5"))

	cnt, _, err = db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condEqual2, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(cnt))
	condEqual5 := &IndexCondition{
		StartKey:     []byte("5"),
		IncludeStart: true,
		EndKey:       []byte("5"),
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}
	cnt, pkList5, err := db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condEqual5, false)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, 1, len(pkList5))
	assert.Equal(t, pkList[0].PKey, pkList5[0].PKey)
	assert.Equal(t, condEqual5.StartKey, pkList5[0].IndexValue)

	cnt, _, err = db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))

	// hincrby
	db.HIncrBy(0, inputPKList[1], intIndex.IndexField, 1)
	cnt, _, err = db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condEqual5, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(cnt))
	condEqual6 := &IndexCondition{
		StartKey:     []byte("6"),
		IncludeStart: true,
		EndKey:       []byte("6"),
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}
	cnt, pkList6, err := db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condEqual6, false)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, 1, len(pkList6))
	assert.Equal(t, pkList[0].PKey, pkList6[0].PKey)
	assert.Equal(t, condEqual6.StartKey, pkList6[0].IndexValue)

	// hmset test
	// hdel
	// hclear

	inputFVList := make([][]byte, 0, 3)
	inputFVList = append(inputFVList, []byte("fv1"))
	inputFVList = append(inputFVList, []byte("fv2"))
	inputFVList = append(inputFVList, []byte("fv3"))
	db.wb.Clear()
	for i, pk := range inputPKList {
		err = db.HMset(0, pk, common.KVRecord{stringIndex.IndexField, inputFVList[i]})
		assert.Nil(t, err)
	}

	cnt, _, err = db.HsetIndexSearch(stringIndex.Table, stringIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList), int(cnt))

	condEqual0 := &IndexCondition{
		StartKey:     inputFVList[0],
		IncludeStart: true,
		EndKey:       inputFVList[0],
		IncludeEnd:   true,
		Offset:       0,
		Limit:        -1,
	}

	cnt, pkList, err = db.HsetIndexSearch(stringIndex.Table, stringIndex.IndexField, condEqual0, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, int(cnt))
	assert.Equal(t, 1, len(pkList))
	assert.Equal(t, inputPKList[0], pkList[0].PKey)

	db.HDel(inputPKList[0], stringIndex.IndexField)
	cnt, _, err = db.HsetIndexSearch(stringIndex.Table, stringIndex.IndexField, condEqual0, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(cnt))

	cnt, _, err = db.HsetIndexSearch(stringIndex.Table, stringIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))

	db.HClear(inputPKList[1])

	cnt, _, err = db.HsetIndexSearch(intIndex.Table, intIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-1, int(cnt))

	cnt, _, err = db.HsetIndexSearch(stringIndex.Table, stringIndex.IndexField, condAll, false)
	assert.Nil(t, err)
	assert.Equal(t, len(inputPKList)-2, int(cnt))
}
