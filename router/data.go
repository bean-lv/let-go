package router

const (
	Static_Folder = "www"
	Project_Name  = "my-zone"
	Homepage      = "index.html"

	Prefix_API    = "/api"
	Prefix_Static = "/www"
	Prefix_Upload = "/upload"

	Suffix_Controller = "Controller"
)

var (
	HTTPMethods = []string{
		"CONNECT",
		"DELETE",
		"GET",
		"HEAD",
		"OPTIONS",
		"PATCH",
		"POST",
		"PROPFIND",
		"PUT",
		"TRACE",
	}
)
