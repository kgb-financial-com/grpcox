package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kgb-financial-com/grpcox/core"
)

// Handler hold all handler methods
type Handler struct {
	g *core.GrpCox
}

// InitHandler Constructor
func InitHandler() *Handler {
	return &Handler{
		g: core.InitGrpCox(),
	}
}

func (h *Handler) getActiveConns(w http.ResponseWriter, r *http.Request) {
	response(w, h.g.GetActiveConns())
}

func (h *Handler) closeActiveConns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	err := h.g.CloseActiveConns(strings.Trim(host, " "))
	if err != nil {
		writeError(w, err)
		return
	}
	response(w, map[string]bool{"success": true})
}

func (h *Handler) getKnownServers(w http.ResponseWriter, r *http.Request) {
	response(w, h.g.GetKnownServers())
}

func (h *Handler) getLists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	service := vars["serv_name"]

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))
	restart, _ := strconv.ParseBool(r.FormValue("restart"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, restart)
	if err != nil {
		writeError(w, err)
		return
	}

	result, err := res.List(service)
	if err != nil {
		writeError(w, err)
		return
	}

	h.g.Extend(host)
	response(w, result)
}

// getListsWithProto handling client request for service list with proto
func (h *Handler) getListsWithProto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	service := vars["serv_name"]

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))
	restart, _ := strconv.ParseBool(r.FormValue("restart"))

	// limit upload file to 5mb
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		writeError(w, err)
		return
	}

	// convert uploaded files to list of Proto struct
	files := r.MultipartForm.File["protos"]
	protos := make([]core.Proto, 0, len(files))
	for _, file := range files {
		fileData, err := file.Open()
		if err != nil {
			writeError(w, err)
			return
		}
		defer fileData.Close()

		content, err := ioutil.ReadAll(fileData)
		if err != nil {
			writeError(w, err)
		}

		protos = append(protos, core.Proto{
			Name:    file.Filename,
			Content: content,
		})
	}

	res, err := h.g.GetResourceWithProto(context.Background(), host, !useTLS, restart, protos)
	if err != nil {
		writeError(w, err)
		return
	}

	result, err := res.List(service)
	if err != nil {
		writeError(w, err)
		return
	}

	h.g.Extend(host)
	response(w, result)
}

func (h *Handler) describeFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	funcName := vars["func_name"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Func Name"))
		return
	}

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		writeError(w, err)
		return
	}

	// get param
	result, _, err := res.Describe(funcName)
	if err != nil {
		writeError(w, err)
		return
	}
	match := reGetFuncArg.FindStringSubmatch(result)
	if len(match) < 2 {
		writeError(w, fmt.Errorf("Invalid Func Type"))
		return
	}

	// describe func
	result, template, err := res.Describe(match[1])
	if err != nil {
		writeError(w, err)
		return
	}

	type desc struct {
		Schema   string `json:"schema"`
		Template string `json:"template"`
	}

	h.g.Extend(host)
	response(w, desc{
		Schema:   result,
		Template: template,
	})

}

func (h *Handler) invokeFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Host"))
		return
	}

	funcName := vars["func_name"]
	if host == "" {
		writeError(w, fmt.Errorf("Invalid Func Name"))
		return
	}

	useTLS, _ := strconv.ParseBool(r.Header.Get("use_tls"))

	res, err := h.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		writeError(w, err)
		return
	}

	// context metadata
	metadataHeader := r.Header.Get("Metadata")
	metadataArr := strings.Split(metadataHeader, ",")

	// construct array of string with "key: value" form to satisfy grpcurl MetadataFromHeaders
	var metadata []string
	var metadataStr string
	for i, m := range metadataArr {
		i += 1
		if isEven := i%2 == 0; isEven {
			metadataStr = metadataStr + m
			metadata = append(metadata, metadataStr)
			metadataStr = ""
			continue
		}
		metadataStr = fmt.Sprintf("%s:", m)
	}

	// get param
	result, timer, err := res.Invoke(context.Background(), metadata, funcName, r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	type invRes struct {
		Time   string `json:"timer"`
		Result string `json:"result"`
	}

	h.g.Extend(host)
	response(w, invRes{
		Time:   timer.String(),
		Result: result,
	})
}
