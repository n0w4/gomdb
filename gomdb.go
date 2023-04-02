package gomdb

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/google/uuid"
)

type document map[string]interface{}

type MemoryDB struct {
	collections map[string][]document
	Name        string
	mu          sync.RWMutex
}

func NewMemoryDB(dbName string) *MemoryDB {
	collections := make(map[string][]document)

	return &MemoryDB{
		collections: collections,
		Name:        dbName,
	}
}

func (mdb *MemoryDB) InsertOnCollection(collectionName string, doc map[string]interface{}) {
	doc = mdb.injectId(doc)
	doc = mdb.injectFields(doc)
	doc = mdb.injectInSync(doc)

	mdb.collections[collectionName] = mdb.getCollection(collectionName)
	mdb.collections[collectionName] = append(mdb.collections[collectionName], doc)
}

func (mdb *MemoryDB) getCollection(collectionName string) []document {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()
	if ok := mdb.collections[collectionName]; ok == nil {
		mdb.collections[collectionName] = make([]document, 0)
	}
	return mdb.collections[collectionName]
}

func (mdb *MemoryDB) injectId(doc map[string]interface{}) map[string]interface{} {
	doc["_id"] = uuid.New().String()
	return doc
}

func (mdb *MemoryDB) injectFields(doc map[string]interface{}) map[string]interface{} {
	fields := make([]string, 0)

	for k := range doc {
		fields = append(fields, k)
	}

	doc["_fields"] = fields
	return doc
}

func (mdb *MemoryDB) injectInSync(doc map[string]interface{}) map[string]interface{} {
	if doc["_in_sync"] != nil {
		return doc
	}

	doc["_in_sync"] = false
	return doc
}

func (mdb *MemoryDB) FindOnCollection(collectionName string, filter map[string]interface{}) []document {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()

	collection := mdb.getCollection(collectionName)
	filteredDocs := make([]document, 0)

	for _, doc := range collection {
		if mdb.parseDocument(doc, filter) {
			filteredDocs = append(filteredDocs, doc)
		}
	}

	return filteredDocs
}

func (mdb *MemoryDB) UpdateOnCollection(collectionName string, filter, update map[string]interface{}) int {
	collection := mdb.getCollection(collectionName)
	updatedCount := 0

	for i, doc := range collection {
		if mdb.parseDocument(doc, filter) {
			mdb.makeChange(&doc, update)
			collection[i] = doc
			updatedCount++
		}
	}

	return updatedCount
}

func (mdb *MemoryDB) makeChange(doc *document, update map[string]interface{}) {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	for key, value := range update {
		if !mdb.keyCanBeChanged(key) {
			continue
		}
		(*doc)[key] = value
	}
}

func (mdb *MemoryDB) keyCanBeChanged(key string) bool {
	return key != "_id" && key != "_fields" && key != "_in_sync"
}

func (mdb *MemoryDB) parseDocument(doc document, filter map[string]interface{}) bool {
	matchs := 0

	for filterKey, filterValue := range filter {
		docValue, ok := doc[filterKey]
		if !ok {
			continue
		}

		switch filterValue.(type) {
		case string:
			rx := regexp.MustCompile(filterValue.(string))
			if !rx.MatchString(docValue.(string)) {
				continue
			}
		default:
			if !reflect.DeepEqual(docValue, filterValue) {
				continue
			}
		}

		matchs++
	}

	return matchs == len(filter)
}
