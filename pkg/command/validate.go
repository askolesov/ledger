package command

import (
	"go-finances/pkg/model_v1"

	"github.com/spf13/cobra"
)

func getValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate your finances data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			data, err := model_v1.ReadData(path)
			if err != nil {
				return err
			}

			cmd.Println(data)

			err = data.Validate()
			if err != nil {
				return err
			}

			cmd.Println("Data is valid")

			return nil
		},
	}
}
