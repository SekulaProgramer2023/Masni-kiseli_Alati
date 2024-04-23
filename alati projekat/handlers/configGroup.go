package handlers

import (
	"alati_projekat/model"
	"alati_projekat/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"strconv"
)

type ConfigGroupHandler struct {
	service       services.ConfigGroupService
	serviceConfig services.ConfigService
}

func NewConfigGruopHandler(service services.ConfigGroupService, serviceConfig services.ConfigService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service:       service,
		serviceConfig: serviceConfig,
	}
}

func decodeBodyCG(r io.Reader) (*model.ConfigGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt model.ConfigGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionInt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content−Type", "application/json")
	w.Write(resp)
}

func (c ConfigGroupHandler) Add(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyCG(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.service.Add(*rt)

	renderJSON(w, rt)
}

func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(name, version)
	if err != nil {
		http.Error(w, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	renderJSON(w, "Deleted")
}

func (c ConfigGroupHandler) AddConfToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	conf, err := c.serviceConfig.Get(nameC, versionC)
	if err != nil {
		http.Error(w, "Failed to fetch config", http.StatusInternalServerError)
		return
	}

	group, err := c.service.Get(nameG, versionG)
	if err != nil {
		http.Error(w, "Failed to fetch config group", http.StatusInternalServerError)
		return
	}

	group.Configs = append(group.Configs, conf)
	c.service.Add(group)

	renderJSON(w, "success Put")
}

func (c ConfigGroupHandler) RemoveConfFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	group, err := c.service.Get(nameG, versionG)
	if err != nil {
		http.Error(w, "Failed to fetch config group", http.StatusInternalServerError)
		return
	}

	var updatedConfigs []model.Config
	for _, conf := range group.Configs {
		confKey := fmt.Sprintf("%s/%d", conf.Name, conf.Version)
		key := fmt.Sprintf("%s/%d", nameC, versionC)
		if confKey != key {
			updatedConfigs = append(updatedConfigs, conf)
		}
	}
	group.Configs = updatedConfigs

	c.service.Add(group)

	renderJSON(w, "success Put")
}

func (c ConfigGroupHandler) renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
