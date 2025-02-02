package database

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	postgrest "github.com/supabase-community/postgrest-go"
	"log"
	"maps"
	s "strings"
)

type Export struct {
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var p = fmt.Printf

func (d Export) GetUserID() string {
	db, err := badger.Open(badger.DefaultOptions("./tmp/badger"))
	handleErr(err)

	id := []byte{102}

	verr := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ID"))
		if err != nil {
			p("The result of ID: %s\n", err)
			db.Update(func(txn *badger.Txn) error {
				newID := uuid.New().String()
				e := badger.NewEntry([]byte("ID"), []byte(newID))
				return txn.SetEntry(e)
			})
		} else {
			item.Value(func(val []byte) error {
				id = append([]byte{}, val...)
				return nil
			})
		}
		return nil
	})
	handleErr(verr)

	defer db.Close()
	return string(id)
}

const auth string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiZm9yZWlnbl91c2VyIn0.xYNvBmyr6PxqXayZQlIOHw5zk1pi33njxDvbFvoesRg"

func (d Export) InitDB(ip, port, schema string) *postgrest.Client {

	c := postgrest.NewClient("http://"+ip+":"+port, schema, map[string]string{
		"Authorization": "Bearer " + auth,
	})
	if c.ClientError != nil {
		panic(c.ClientError)
	}
	return c
}

func (d Export) CheckPostDB(c *postgrest.Client, str string) bool {

	mapstr := stringtokv(str)
	res, _, err := c.
		From("orders").
		Select("uuid, link, description, delivery_address, crypto", "", false).
		Eq("uuid", mapstr["uuid"]).
		Eq("link", mapstr["link"]).
		Eq("description", mapstr["description"]).
		Eq("delivery_address", mapstr["delivery_address"]).
		Eq("crypto", mapstr["crypto"]).
		ExecuteString()

	runeRes := []rune(res)

	if err != nil {
		p("error from clientPostgrest: %s\n", err)
		return true
	}
	if len(runeRes) < 4 || res == "" {
		p("no result from clientPostgrest")
		return true
	}
	return false
}

func (d Export) CheckUpdateTxidDB(c *postgrest.Client, mtx []string) bool {
	res, _, err := c.
		From("orders").
		Select("uuid, txid", "", false).
		Eq("uuid", mtx[1]).
		Eq("txid", mtx[2]).
		ExecuteString()

	runeRes := []rune(res)

	if err != nil {
		p("error from clientPostgrest: %s\n", err)
		return false
	}
	if len(runeRes) < 4 || res == "" {
		p("no result from clientPostgrest")
		return false
	}
	return true
}

func (d Export) SelectDB(c *postgrest.Client, table, userID string, keys []string) string {

	var skeys string = s.Join(keys, ", ")

	res, _, err := c.
		From(table).
		Select(skeys, "", false).
		Eq("uuid", userID).
		ExecuteString()

	runeRes := []rune(res)

	if err != nil {
		p("error from clientPostgrest: %s\n", err)
		return ""
	}
	if len(runeRes) < 4 || res == "" {
		p("no result from clientPostgrest")
		return ""
	}
	return res
}

func (d Export) UpdateDB(c *postgrest.Client, table, userID, id string, toUpdate map[string]string) ([]byte, error) {
	res, _, err := c.
		From(table).
		Update(toUpdate, "", "").
		Eq("uuid", userID).
		Eq("id", id).
		Execute()

	return res, err
}

func stringtokv(str string) map[string]string {
	a := s.Split(str, ", ")

	mapstr := make(map[string]string, 6)
	for i := 0; i < len(a); i++ {
		k, v, _ := s.Cut(a[i], ": ")
		data := map[string]string{k: v}
		maps.Copy(mapstr, data)
	}

	return mapstr
}

func (d Export) InsertDB(c *postgrest.Client, str string) {
	mapstr := stringtokv(str)
	c.From("orders").
		Insert(map[string]interface{}{
			"uuid":             mapstr["uuid"],
			"link":             mapstr["link"],
			"description":      mapstr["description"],
			"delivery_address": mapstr["delivery_address"],
			"crypto":           mapstr["crypto"],
			"txid":             mapstr["txid"],
		}, false, "", "representation", "exact").
		Execute()
}
