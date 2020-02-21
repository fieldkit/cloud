// Code generated by goa v3.0.10, DO NOT EDIT.
//
// test HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	testc "github.com/fieldkit/cloud/server/api/gen/http/test/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `test get
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` test get --id 3166109748888346624` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		testFlags = flag.NewFlagSet("test", flag.ContinueOnError)

		testGetFlags  = flag.NewFlagSet("get", flag.ExitOnError)
		testGetIDFlag = testGetFlags.String("id", "REQUIRED", "")
	)
	testFlags.Usage = testUsage
	testGetFlags.Usage = testGetUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "test":
			svcf = testFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "test":
			switch epn {
			case "get":
				epf = testGetFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "test":
			c := testc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get":
				endpoint = c.Get()
				data, err = testc.BuildGetPayload(*testGetIDFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// testUsage displays the usage of the test command and its subcommands.
func testUsage() {
	fmt.Fprintf(os.Stderr, `Service is the test service interface.
Usage:
    %s [globalflags] test COMMAND [flags]

COMMAND:
    get: Get implements get.

Additional help:
    %s test COMMAND --help
`, os.Args[0], os.Args[0])
}
func testGetUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] test get -id INT64

Get implements get.
    -id INT64: 

Example:
    `+os.Args[0]+` test get --id 3166109748888346624
`, os.Args[0])
}
