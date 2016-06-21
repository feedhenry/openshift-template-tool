package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/openshift/origin/pkg/template/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api/latest"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/feedhenry/openshift-template-tool/template"
)

// NewMergeCommand creates a new command to merge OpenShift templates.
func NewMergeCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge BASE_TEMPLATE TEMPLATE_1 TEMPLATE_2 ...",
		Short: "Merge OpenShift templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunMerge(cmd, args, stdin, stdout, stderr)
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("output", "o", "json", "Output format. One of: json|yaml.")
	cmd.Flags().String("output-version", latest.Version, "Output objects with the given version.")
	return cmd
}

// RunMerge merges OpenShift template files from args and prints a merged
// template to stdout.
func RunMerge(cmd *cobra.Command, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	loader := &templateLoader{kapi.Codecs.UniversalDecoder(), runtime.UnstructuredJSONScheme}
	templates, err := loader.FromFiles(args...)
	if err != nil {
		return err
	}

	merged := template.Merge(templates...)
	if merged == nil {
		return nil
	}

	version, err := util.OutputVersion(cmd, &latest.ExternalVersion)
	if err != nil {
		return err
	}

	// Explicitly convert template objects to output version, because the
	// printer won't do it.
	for i, obj := range merged.Objects {
		var converted runtime.Object
		converted, err = kapi.Scheme.ConvertToVersion(obj, version.String())
		if err != nil {
			return err
		}
		merged.Objects[i] = converted
	}

	outputFormat := util.GetFlagString(cmd, "output")
	if outputFormat != "json" && outputFormat != "yaml" {
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	p, _, err := kubectl.GetPrinter(outputFormat, "")
	if err != nil {
		return err
	}
	p = kubectl.NewVersionedPrinter(p, kapi.Scheme, version)
	return p.PrintObj(merged, stdout)
}

// templateLoader loads OpenShift templates from files or bytes.
type templateLoader struct {
	decoder             runtime.Decoder
	unstructuredDecoder runtime.Decoder
}

func (tl *templateLoader) FromFiles(paths ...string) ([]*api.Template, error) {
	var templates []*api.Template
	for _, path := range paths {
		tmpl, err := tl.FromFile(path)
		if err != nil {
			return nil, err
		}
		templates = append(templates, tmpl)
	}
	return templates, nil
}

func (tl *templateLoader) FromFile(path string) (*api.Template, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tmpl, err := tl.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %v", path, err)
	}
	return tmpl, nil
}

func (tl *templateLoader) Decode(data []byte) (*api.Template, error) {
	dec := tl.decoder

	obj, _, err := dec.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}

	tmpl, ok := obj.(*api.Template)
	if !ok {
		kind := reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
		return nil, fmt.Errorf("top level object must be of kind Template, found %s", kind)
	}

	return tmpl, tl.resolveObjects(tmpl)
}

func (tl *templateLoader) resolveObjects(tmpl *api.Template) error {
	dec := tl.decoder
	udec := tl.unstructuredDecoder
	for i, obj := range tmpl.Objects {
		if unknown, ok := obj.(*runtime.Unknown); ok {
			decoded, _, err := dec.Decode(unknown.Raw, nil, nil)
			if err != nil {
				debugf("ignoring API type checking of %s because of error: %v", unknown.Raw, err)
				decoded, _, err = udec.Decode(unknown.Raw, nil, nil)
				if err != nil {
					return err
				}
			}
			tmpl.Objects[i] = decoded
		}
	}
	return nil
}

// debugf prints messages to stderr if the program is run with a debug flag.
func debugf(format string, a ...interface{}) {
	if d := pflag.Lookup("debug"); d != nil && d.Value.String() == "true" {
		if format[len(format)-1] != '\n' {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format, a...)
	}
}
