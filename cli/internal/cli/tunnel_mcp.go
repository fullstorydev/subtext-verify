// tunnel_mcp.go has no build constraint: the in-process tunnel engine works
// on all platforms (no daemon/POSIX-signal machinery needed here).

package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/fullstorydev/subtext-verify/cli/internal/tunnelmcp"
)

var tunnelMcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start an MCP stdio server exposing tunnel-connect, tunnel-disconnect, and tunnel-status",
	Long: `Start a Model Context Protocol server over stdio that exposes three tools:

  tunnel-connect     Connect a local dev server to the Subtext relay
  tunnel-disconnect  Disconnect one or all active tunnels
  tunnel-status      List all active tunnels

This command is the stdio server backing the subtext-tunnel MCP server entry
in .mcp.json. It is launched automatically by Claude Code (and other MCP hosts)
via npx:

  npx -y @subtextdev/subtext-cli@latest tunnel mcp

No API key is required. The relay URL supplied to tunnel-connect already carries
auth (embedded token from live-tunnel). All tunnel logging goes to stderr;
stdout is reserved for JSON-RPC.`,
	RunE: runTunnelMcp,
}

func runTunnelMcp(cmd *cobra.Command, _ []string) error {
	srv := tunnelmcp.New(binaryVersion())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	return srv.Serve(ctx, os.Stdin, os.Stdout)
}
