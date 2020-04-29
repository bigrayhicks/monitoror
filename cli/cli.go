package cli

import (
	"io"
	"os"

	"github.com/monitoror/monitoror/store"

	"github.com/spf13/cobra"
)

type MonitororCli struct {
	rootCmd *cobra.Command
	store   *store.Store
	output  io.Writer
}

func NewMonitororCli(store *store.Store) *MonitororCli {
	return &MonitororCli{
		store:  store,
		output: os.Stdout,
	}
}

func (c *MonitororCli) SetRootCommand(rootCmd *cobra.Command) {
	c.rootCmd = rootCmd
}

func (c *MonitororCli) GetRootCommand() *cobra.Command {
	return c.rootCmd
}

func (c *MonitororCli) GetStore() *store.Store {
	return c.store
}

func (c *MonitororCli) GetOutput() io.Writer {
	return c.output
}
