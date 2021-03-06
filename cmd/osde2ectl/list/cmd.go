package list

import (
	"fmt"

	"github.com/openshift/osde2e/cmd/osde2e/common"
	"github.com/openshift/osde2e/pkg/common/clusterproperties"
	"github.com/openshift/osde2e/pkg/common/metadata"
	"github.com/openshift/osde2e/pkg/common/providers"
	"github.com/openshift/osde2e/pkg/common/spi"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List created/existing clusters made by osde2e",
	Long:  "List specific clusters using the provided arguments.",
	Args:  cobra.OnlyValidArgs,
	RunE:  run,
}

var args struct {
	configString    string
	customConfig    string
	secretLocations string
}

func init() {
	flags := Cmd.Flags()

	flags.StringVar(
		&args.configString,
		"configs",
		"",
		"A comma separated list of built in configs to use",
	)
	flags.StringVar(
		&args.customConfig,
		"custom-config",
		"",
		"Custom config file for osde2e",
	)
	flags.StringVar(
		&args.secretLocations,
		"secret-locations",
		"",
		"A comma separated list of possible secret directory locations for loading secret configs.",
	)

	Cmd.RegisterFlagCompletionFunc("output-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "prom"}, cobra.ShellCompDirectiveDefault
	})
}

func run(cmd *cobra.Command, argv []string) error {
	var provider spi.Provider
	var err error
	if err = common.LoadConfigs(args.configString, args.customConfig, args.secretLocations); err != nil {
		return fmt.Errorf("error loading initial state: %v", err)
	}

	if provider, err = providers.ClusterProvider(); err != nil {
		return fmt.Errorf("could not setup cluster provider: %v", err)
	}

	metadata.Instance.SetEnvironment(provider.Environment())

	clusters, err := provider.ListClusters("properties.MadeByOSDe2e='true'")
	if err != nil {
		return err
	}

	fmt.Printf("%-25s%-35s%-15s%-25s%-25s%-25s%s\n", "NAME", "ID", "STATE", "STATUS", "OWNER", "INSTALLED VERSION", "UPGRADE VERSION")
	for _, cluster := range clusters {
		properties := cluster.Properties()
		fmt.Printf("%-25s%-35s%-15s%-25s%-25s%-25s%s\n", cluster.Name(), cluster.ID(), cluster.State(), properties[clusterproperties.Status], properties[clusterproperties.OwnedBy], properties[clusterproperties.InstalledVersion], properties[clusterproperties.UpgradeVersion])
	}

	return nil
}
