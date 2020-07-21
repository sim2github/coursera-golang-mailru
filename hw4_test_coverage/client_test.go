package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
)

type Person struct {
	User
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

type Users []*Person

func (u Users) Len() int      { return len(u) }
func (u Users) Swap(i, j int) { u[i], u[j] = u[j], u[i] }

type ByName struct{ Users }

func (u ByName) Less(i, j int) bool { return u.Users[i].Name < u.Users[j].Name }

type ByAge struct{ Users }

func (u ByAge) Less(i, j int) bool { return u.Users[i].Age < u.Users[j].Age }

var (
	users Users
	ts    *httptest.Server
)

func init() {
	file, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		log.Fatalf("could not read file: %v", err)
	}
	data := xml.NewDecoder(bytes.NewReader(file))

	for {
		t, _ := data.Token()
		if t == nil {
			break
		}
		switch el := t.(type) {
		case xml.StartElement:
			if el.Name.Local == "row" {
				var row Person
				data.DecodeElement(&row, &el)
				row.Name = row.FirstName + " " + row.LastName
				users = append(users, &row)
			}
		default:
		}
	}
	ts = httptest.NewServer(http.HandlerFunc(SearchServer))
}

func parseParams(r *http.Request) SearchRequest {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit < 0 {
		limit = 0
	}
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	orderBy, err := strconv.Atoi(r.FormValue("order_by"))
	if err != nil || (orderBy != -1 && orderBy != 1) {
		orderBy = 0
	}

	orderField := r.FormValue("order_field")
	if orderField == "" || (orderField != "Id" && orderField != "Age") {
		orderField = "Name"
	}

	query := r.FormValue("query")

	return SearchRequest{
		Limit:      limit,
		Offset:     offset,
		Query:      query,
		OrderField: orderField,
		OrderBy:    orderBy,
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Header["Accesstoken"][0] == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(SearchErrorResponse{Error: "Acces denied"})
		return
	}

	req := parseParams(r)

	if req.OrderField == "Name" {
		sort.Sort(ByName{users})
	} else if req.OrderField == "Age" {
		sort.Sort(ByAge{users})
	}

	// else {
	// 	sort.SliceStable(users, func(i, j int) bool {
	// 		return users[i].Id < users[j].Id
	// 	})
	// }

	u := users[:req.Limit]
	// w.WriteHeader(http.StatusInternalServerError)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func TestFindUsers(t *testing.T) {
	defer ts.Close()
	var (
		ac = SearchClient{
			URL:         ts.URL,
			AccessToken: "#admin_access_token",
		}

		c = SearchClient{
			URL: ts.URL,
		}

		tests = []struct {
			name string
			cl   *SearchClient
			req  SearchRequest
			rsp  *SearchResponse
			err  string
		}{
			{"Access token", &c, SearchRequest{Limit: 1}, &SearchResponse{}, "Bad AccessToken"},
			{"Default list", &ac, SearchRequest{OrderField: "Name", Limit: 1}, &SearchResponse{Users: []User{users[0].User}, NextPage: true}, ""},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp, err := tt.cl.FindUsers(tt.req)
			if err != nil {
				if err.Error() != tt.err {
					t.Errorf("(%+v): expected error '%s', get '%s'", tt.req, tt.err, err.Error())
					return
				}
				return
			} else {
				if rsp != tt.rsp {
					fmt.Printf("users[0] = %+v\n", users[0])
					t.Errorf("(%+v): expected response '%+v', get '%+v'", tt.req, tt.rsp, rsp)
				}
			}
		})
	}
}
