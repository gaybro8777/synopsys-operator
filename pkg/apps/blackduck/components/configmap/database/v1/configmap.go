package v1

import (
	"fmt"
	"strconv"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/store"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/types"
	apputils "github.com/blackducksoftware/synopsys-operator/pkg/apps/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"k8s.io/client-go/kubernetes"
)

// BdConfigmap holds the Black Duck config map configuration
type BdConfigmap struct {
	config     *protoform.Config
	kubeClient *kubernetes.Clientset
	blackDuck  *blackduckapi.Blackduck
}

// NewBdConfigmap returns the Black Duck config map configuration
func NewBdConfigmap(config *protoform.Config, kubeClient *kubernetes.Clientset, cr interface{}) (types.ConfigMapInterface, error) {
	blackDuck, ok := cr.(*blackduckapi.Blackduck)
	if !ok {
		return nil, fmt.Errorf("unable to cast the interface to Black Duck object")
	}
	return &BdConfigmap{config: config, kubeClient: kubeClient, blackDuck: blackDuck}, nil
}

// GetCM returns the config map
func (b BdConfigmap) GetCM() []*components.ConfigMap {
	var configMaps []*components.ConfigMap
	// DB
	hubDbConfig := components.NewConfigMap(horizonapi.ConfigMapConfig{Namespace: b.blackDuck.Spec.Namespace, Name: apputils.GetResourceName(b.blackDuck.Name, util.BlackDuckName, "db-config")})
	if b.blackDuck.Spec.ExternalPostgres != nil {
		hubDbConfig.AddData(map[string]string{
			"HUB_POSTGRES_ADMIN": b.blackDuck.Spec.ExternalPostgres.PostgresAdmin,
			"HUB_POSTGRES_USER":  b.blackDuck.Spec.ExternalPostgres.PostgresUser,
			"HUB_POSTGRES_PORT":  strconv.Itoa(b.blackDuck.Spec.ExternalPostgres.PostgresPort),
			"HUB_POSTGRES_HOST":  b.blackDuck.Spec.ExternalPostgres.PostgresHost,
		})
	} else {
		hubDbConfig.AddData(map[string]string{
			"HUB_POSTGRES_ADMIN": "blackduck",
			"HUB_POSTGRES_USER":  "blackduck_user",
			"HUB_POSTGRES_PORT":  "5432",
			"HUB_POSTGRES_HOST":  apputils.GetResourceName(b.blackDuck.Name, util.BlackDuckName, "postgres"),
		})
	}

	if b.blackDuck.Spec.ExternalPostgres != nil {
		hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL": strconv.FormatBool(b.blackDuck.Spec.ExternalPostgres.PostgresSsl)})
		if b.blackDuck.Spec.ExternalPostgres.PostgresSsl {
			hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL_CERT_AUTH": "false"})
		}
	} else {
		hubDbConfig.AddData(map[string]string{"HUB_POSTGRES_ENABLE_SSL": "false"})
	}
	hubDbConfig.AddLabels(apputils.GetVersionLabel("postgres", b.blackDuck.Name, b.blackDuck.Spec.Version))
	configMaps = append(configMaps, hubDbConfig)
	return configMaps
}
func init() {
	store.Register(types.BlackDuckDatabaseConfigmapV1, NewBdConfigmap)
}
