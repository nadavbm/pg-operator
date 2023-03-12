package specs

import (
	"fmt"
	"strings"

	configmaps "example.com/pg/apis/configmaps/v1alpha1"
)

func generateHbaConf(cm *configmaps.ConfigMap) string {
	config := ""
	for _, v := range cm.Spec.HbaConf {
		config = config + fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", v.Type, v.Database, v.User, v.Address, v.IpMask, v.Method)
	}

	return config
}

func generatePostgresqlConf(cm *configmaps.ConfigMap) string {
	// postgresql.conf server settings - sample https://github.com/postgres/postgres/blob/master/src/backend/utils/misc/postgresql.conf.sample
	config := fmt.Sprintf("data_directory = %s\n", cm.Spec.PGConf.DataDir)
	config += fmt.Sprintf("hba_file = %s\n", cm.Spec.PGConf.HbaFile)
	config += fmt.Sprintf("max_connections = %d\n", cm.Spec.PGConf.MaxConnections)
	config += fmt.Sprintf("port = %d\n", cm.Spec.PGConf.Port)
	cfg := strings.Replace(config, `\n`, "\n", -1)
	return cfg
}
