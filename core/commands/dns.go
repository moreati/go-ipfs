package commands

import (
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	namesys "github.com/ipfs/go-ipfs/namesys"
	util "github.com/ipfs/go-ipfs/util"
)

var DNSCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "DNS link resolver",
		ShortDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, IPFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.
`,
		LongDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, IPFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.

For example, with this DNS TXT record:

  ipfs.io. TXT "dnslink=/ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy ..."

The resolver will give:

  > ipfs dns ipfs.io
  /ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy

And with this DNS TXT record:

  ipfs.ipfs.io. TXT "dnslink=/dns/ipfs.io ..."

The resolver will give:

  > ipfs dns ipfs.io
  /dns/ipfs.io
  > ipfs dns --recursive
  /ipfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("domain-name", true, false, "The domain-name name to resolve.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption("recursive", "r", "Resolve until the result is not a DNS link"),
	},
	Run: func(req cmds.Request, res cmds.Response) {

		recursive, _, _ := req.Option("recursive").Bool()
		name := req.Arguments()[0]
		resolver := namesys.NewDNSResolver()

		depth := 1
		if recursive {
			depth = namesys.DefaultDepthLimit
		}
		output, err := resolver.ResolveN(req.Context(), name, depth)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}
		res.SetOutput(&ResolvedPath{output})
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			output, ok := res.Output().(*ResolvedPath)
			if !ok {
				return nil, util.ErrCast()
			}
			return strings.NewReader(output.Path.String()), nil
		},
	},
	Type: ResolvedPath{},
}
