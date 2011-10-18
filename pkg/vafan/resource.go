package vafan

func get(r request) (output string) {

    // get this resource's data
    d = r.resource.getData();

    // put the data through the templates, and return...
    output = parsePath

    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */
    return output
}

type resource interface {
    getData() schema
}

type schema struct {

}

// video resource
type video struct {
    parts []string
}

func (r *video) getData() (s schema) {

    // get the relevant data 

}

