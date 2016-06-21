package template

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"

	"github.com/openshift/origin/pkg/template/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/runtime"
)

var debug = false

// Merge merges the parameters and objects of input templates into the first
// template and returns the modified first template.
func Merge(templates ...*api.Template) *api.Template {
	if len(templates) == 0 {
		return nil
	}
	base := templates[0]
	for _, tmpl := range templates[1:] {
		base.Parameters = append(base.Parameters, tmpl.Parameters...)
		base.Objects = append(base.Objects, tmpl.Objects...)
	}
	// Remove duplicates from Parameters.
	{
		seen := make(map[string]struct{}, len(base.Parameters))
		var key string
		for _, p := range base.Parameters {
			key = p.Name
			if _, ok := seen[key]; !ok {
				base.Parameters[len(seen)] = p
				seen[key] = struct{}{}
			}
		}
		base.Parameters = base.Parameters[:len(seen)]
	}
	// Remove duplicates from Objects.
	{
		seen := make(map[string]struct{}, len(base.Objects))
		var key string
		for _, o := range base.Objects {
			switch typed := o.(type) {
			case *runtime.Unstructured:
				key = fmt.Sprintf("%s\x00%s", typed.GetKind(), typed.GetName())
			case meta.Object:
				kind := reflect.Indirect(reflect.ValueOf(typed)).Type().Name()
				key = fmt.Sprintf("%s\x00%s", kind, typed.GetName())
			default:
				kind := reflect.Indirect(reflect.ValueOf(typed)).Type().Name()
				key = fmt.Sprintf("%s\x00%d", kind, rand.Int())
			}
			if _, ok := seen[key]; !ok {
				base.Objects[len(seen)] = o
				seen[key] = struct{}{}
			}
		}
		base.Objects = base.Objects[:len(seen)]

		if debug {
			var s []string
			for k := range seen {
				s = append(s, k)
			}
			fmt.Fprintf(os.Stderr, "[DEBUG] objects seen: %#v\n", s)
		}
	}
	return base
}
