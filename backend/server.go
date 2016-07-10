package backend

import (
	"log"

	"strconv"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

//API ..
type API struct{}

//Item ..
type Item struct {
	Key       *datastore.Key `json:"id" datastore:"-"`
	Name      string         `json:"name"`
	Qty       int32          `json:"qty"`
	QtyType   string         `json:"qty_type"`
	Completed bool           `json:"completed"`
	AddedOn time.Time	`json:"added_on"`
}

//Items ..
type ItemsList struct{
	Items []Item `json:"items"`
}

//GetItemReq ..
type GetItemReq struct {
	Key *datastore.Key `json:"id"`
}

type AddItemReq struct {
	Name    string `json:"name"`
	Qty     int32  `json:"qty,string"`
	QtyType string `json:"qty_type"`
}

type UpdateItemReq struct {
	Key       *datastore.Key `json:"id"`
	Completed string         `json:"completed"`
}

type ClearedList struct{
	ClearedItems []Item `json:"cleared_items"`
}


// GetItems ...
func (API) GetItems(c context.Context) (*ItemsList, error) {

	items := []Item{}

	keys, err := datastore.NewQuery("Item").Order("-AddedOn").GetAll(c, &items)

	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		items[i].Key = k
	}

	return &ItemsList{items}, nil

}


func (API) AddItem(c context.Context, r *AddItemReq) (*Item, error) {

	k := datastore.NewIncompleteKey(c, "Item", nil)

	item := Item{Name: r.Name, Qty: r.Qty, QtyType: r.QtyType, AddedOn:time.Now()}

	k, err := datastore.Put(c, k, &item)

	if err != nil {
		return nil, err
	}

	item.Key = k

	return &item, nil
}





func (API) UpdateItem(c context.Context, r *UpdateItemReq) (*Item, error) {
	item, key := Item{}, r.Key

	if err := datastore.Get(c, key, &item); err != nil {
		return nil, err
	}

	item.Completed, _ = strconv.ParseBool(r.Completed)

	key, err := datastore.Put(c, key, &item)

	if err != nil {
		return nil, err
	}

	item.Key = r.Key

	return &item, nil
}



func (API) ClearCompleted(c context.Context)(*ClearedList, error){

	items := []Item{}

	keys, err := datastore.NewQuery("Item").Filter("Completed =", true).GetAll(c, &items);

	if err != nil {
		return nil, err
	}

	 if err := datastore.DeleteMulti(c, keys); err != nil{
	 	return nil, err
	 }

	for i, k := range keys {
		items[i].Key = k
	}


	return &ClearedList{items}, nil;
}



func init() {
	api, err := endpoints.RegisterService(API{}, "sl", "v1", "shopping items", true)

	if err != nil {
		log.Println(err)
	}

	//adapt the name, method, and path for each method.
	info := api.MethodByName("GetItems").Info()
	info.HTTPMethod, info.Path = "GET", "items"


	info = api.MethodByName("AddItem").Info()
	info.HTTPMethod, info.Path = "POST", "items"

	info = api.MethodByName("ClearCompleted").Info()
	info.HTTPMethod, info.Path = "DELETE", "items"

	info = api.MethodByName("UpdateItem").Info()
	info.HTTPMethod, info.Path = "PUT", "items/{id}"

	endpoints.HandleHTTP()
}
