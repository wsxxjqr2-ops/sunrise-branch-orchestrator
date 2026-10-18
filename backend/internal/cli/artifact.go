package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

type artifactReadOptions struct {
	session string
	path    string
}

func newArtifactCommand(ctx *commandContext) *cobra.Command {
	cmd := &cobra.Command{Use: "artifact", Short: "Read artifacts owned by an AO session"}
	var opts artifactReadOptions
	read := &cobra.Command{
		Use:   "read",
		Short: "Read one session artifact as raw text",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ctx.readArtifact(cmd.Context(), opts)
		},
	}
	read.Flags().StringVar(&opts.session, "session", "", "Source session id (required)")
	read.Flags().StringVar(&opts.path, "path", "", "Artifact path relative to the source session (required)")
	cmd.AddCommand(read)
	return cmd
}

func (c *commandContext) readArtifact(ctx context.Context, opts artifactReadOptions) error {
	sessionID := strings.TrimSpace(opts.session)
	if sessionID == "" {
		return usageError{errors.New("usage: ao artifact read --session <session-id> --path <artifact-path>")}
	}
	relative := strings.TrimSpace(opts.path)
	if relative == "" {
		return usageError{errors.New("usage: --path is required")}
	}
	clean := path.Clean(strings.ReplaceAll(relative, "\\", "/"))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return usageError{errors.New("usage: --path must stay inside the source artifact directory")}
	}
	segments := strings.Split(clean, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	requestPath := "sessions/" + url.PathEscape(sessionID) + "/artifact-files/" + strings.Join(segments, "/") + "?source=1"
	response, err := c.getRaw(ctx, requestPath)
	if err != nil {
		return err
	}
	if _, err := c.deps.Out.Write(response.Body); err != nil {
		return err
	}
	if len(response.Body) == 0 || response.Body[len(response.Body)-1] != '\n' {
		_, _ = fmt.Fprintln(c.deps.Out)
	}
	return nil
}
