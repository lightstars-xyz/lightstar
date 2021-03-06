package http

import (
	"github.com/danieldin95/lightstar/src/http/api"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/service"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Host struct {
}

func (h Host) Router(router *mux.Router) {
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("GET")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("POST")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("PUT")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("DELETE")
	router.PathPrefix("/host/{name}/ext/").HandlerFunc(h.Handle).Methods("GET")
}

func (h Host) Filter(w http.ResponseWriter, r *http.Response, data interface{}) bool {
	req := r.Request
	if data == nil || req == nil || req.Method != "GET" {
		return false
	}
	user := data.(*schema.User)
	libstar.Info("Host.Filter %s %s %s", user.Name, req.Method, req.URL.Path)
	if req.URL.Path == "/api/instance" {
		all := schema.ListInstance{
			Items: make([]schema.Instance, 0, 32),
		}
		if err := libstar.GetJSON(r.Body, &all); err != nil {
			libstar.Warn("Host.Filter %s %s", req.Method, req.URL.Path)
			return false
		}
		list := schema.ListInstance{
			Items: make([]schema.Instance, 0, 32),
		}
		obj := api.Instance{}
		for _, item := range all.Items {
			if obj.HasPermission(user, item.Name) {
				list.Items = append(list.Items, item)
			}
		}
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].Name < list.Items[j].Name
		})
		list.Metadata.Size = len(list.Items)
		list.Metadata.Total = len(list.Items)
		api.ResponseJson(w, list)
		return true
	}
	return false
}

func (h Host) Handle(w http.ResponseWriter, r *http.Request) {
	user, _ := api.GetUser(r)
	name, _ := api.GetArg(r, "name")
	node := service.SERVICE.Zone.Get(name)
	if node == nil {
		http.Error(w, "host not found", http.StatusNotFound)
		return
	}
	libstar.Debug("Host.Handle %s", node)
	r.Header.Del("cookie")
	pri := &libstar.ProxyUrl{
		Proxy: libstar.Proxy{
			Prefix: "/host/" + name,
			Server: node.Url,
			Auth: libstar.Auth{
				Type:     "basic",
				Username: node.Username,
				Password: node.Password,
			},
		},
		Filter: h.Filter,
		Data:   &user,
	}
	pri.Initialize()
	pri.Handler(w, r)
}
