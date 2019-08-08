package azurewrapper

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// NsgConfYML structure for nuage info toward vWAN
type NsgConfYML struct {
	NsgData struct {
		Enterprise string   `yaml:"enterprise"`
		NsgName    string   `yaml:"nsg_name"`
		NsgPort    string   `yaml:"nsg_port"`
		PublicIP   string   `yaml:"public_ip"`
		BgpEnabled bool     `yaml:"bgp_enabled"`
		BgpNsgAsn  int64    `yaml:"bgp_nsg_asn"`
		LanSubnet  []string `yaml:"lan_subnet"`
	} `yaml:"nsg_data"`
}

// GetConf function provide a tool to apply the NSG configuration data from a file
func (c *NsgConfYML) GetConf(nsgFile string) *NsgConfYML {

	yamlFile, err := ioutil.ReadFile(nsgFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
