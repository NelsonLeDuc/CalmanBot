package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/nelsonleduc/calmanbot/handlers/models"
)

func HandleActions(w http.ResponseWriter, r *http.Request) {

	absPath, _ := filepath.Abs("./html/actions.html")
	template, _ := template.New("actions.html").Funcs(funcMap).ParseFiles(absPath)

	actions, _ := models.FetchActions(false)
	sort.Sort(models.ByID(actions))

	template.Execute(w, actions)
}

//Template Functions
var funcMap = template.FuncMap{
	"nonNilInt": nonNilInt,
	"nonNilStr": nonNilStr,
}

//Duplicate logic to handle different types semanticallyt
func nonNilInt(val *int) interface{} {
	if val != nil {
		return *val
	}
	return ""
}

func nonNilStr(val *string) string {
	if val != nil {
		return *val
	}
	return ""
}
