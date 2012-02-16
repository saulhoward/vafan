package vafan

func get(r request) (output string) {
    // get this resource's data
    //d = r.resource.getData()
    // put the data through the templates, and return...

    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */

    //from the resource, get the schema
    s := r.resource.getData()
    output = s.name
    return output
}

type resource interface {
    getData() schema
}

type schema struct {
    name string
}

// video resource
type video struct {
    parts []string
}
func (v video) getData() (s schema) {
    // get the relevant data 
    s.name = "zomg"
    return
}

