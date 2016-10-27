package rook

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rook/rook/pkg/model"
	"github.com/rook/rook/pkg/rook/client"
	"github.com/rook/rook/pkg/util/flags"
	"github.com/spf13/cobra"
)

var (
	newImageName     string
	newImagePoolName string
	newImageSize     uint64
)

var blockCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new block image in the cluster",
}

func init() {
	blockCreateCmd.Flags().StringVar(&newImageName, "name", "", "Name of new block image to create (required)")
	blockCreateCmd.Flags().StringVar(&newImagePoolName, "pool-name", "rbd", "Name of storage pool to create new block image in")
	blockCreateCmd.Flags().Uint64Var(&newImageSize, "size", 0, "Size in bytes of the new block image to create (required)")

	blockCreateCmd.MarkFlagRequired("name")
	blockCreateCmd.MarkFlagRequired("size")
	blockCreateCmd.RunE = createBlockImagesEntry
}

func createBlockImagesEntry(cmd *cobra.Command, args []string) error {
	if err := flags.VerifyRequiredFlags(cmd, []string{"name"}); err != nil {
		return err
	}

	if err := flags.VerifyRequiredUint64Flags(cmd, []string{"size"}); err != nil {
		return err
	}

	c := client.NewRookNetworkRestClient(client.GetRestURL(apiServerEndpoint), http.DefaultClient)
	out, err := createBlockImage(newImageName, newImagePoolName, newImageSize, c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(out)
	return nil
}

func createBlockImage(imageName, poolName string, size uint64, c client.RookRestClient) (string, error) {
	newImage := model.BlockImage{Name: imageName, PoolName: poolName, Size: size}
	resp, err := c.CreateBlockImage(newImage)
	if err != nil {
		return "", fmt.Errorf("failed to create new block image '%+v': %+v", newImage, err)
	}

	return resp, nil
}