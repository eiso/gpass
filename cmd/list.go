package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiso/treeprint"
	"github.com/spf13/cobra"
)

type ListCmd struct{}

func NewListCmd() *ListCmd {
	return &ListCmd{}
}

func (c *ListCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all encrypted accounts.",
		RunE:  c.Execute,
	}

	return cmd
}

func (c *ListCmd) Execute(cmd *cobra.Command, args []string) error {
	if err := InitCheck(); err != nil {
		return err
	}

	r := Cfg.Repository

	if err := r.Load(); err != nil {
		return err
	}

	if !r.BranchExists("gpass") {
		return fmt.Errorf("gpass has not been initialized yet, please run: gpass init")
	}

	tree := treeprint.New()
	b := r.ListBranches()

	if len(b) <= 1 {
		return fmt.Errorf("nothing to list here, no accounts have been added yet")
	}

	for _, branch := range b {
		if branch == "gpass" {
			continue
		}

		branch = branch[:len(branch)-4]

		if len(args) > 0 && !strings.HasPrefix(branch, args[0]){
			continue
		}

		parts := strings.Split(branch, string(os.PathSeparator))

		t := tree.FindByValue(parts[0])

		if t == nil {
			t = tree.AddBranch(parts[0])
		}

		for i := 1; i < len(parts); i++ {
			if i == 1 {
				t = t.AddNode(parts[1])
				continue
			}
			t = t.FindLastNode()
			t = t.AddNode(parts[i])
			t = t.FindLastNode()
		}

	}

	fmt.Println(tree.String())

	return nil
}
