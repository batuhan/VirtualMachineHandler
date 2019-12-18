package helpers

import (
	"encoding/base64"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type User struct {
	Name              string
	SshAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
}

type WriteFile struct {
	Encoding string
	Content  string
	Path     string
}

type Template struct {
	PackageUpdate  bool `yaml:"package_update"`
	PackageUpgrade bool `yaml:"package_upgrade"`
	Users          []User
	Chpasswd       struct {
		List   []string
		Expire bool
	}
	Growpart struct {
		Mode                   string
		Devices                []string
		IgnoreGrowrootDisabled bool `yaml:"ignore_growroot_disabled"`
	}
	WriteFiles  []WriteFile `yaml:"write_files"`
	Runcmd      []string
	SshPwauth   bool `yaml:"ssh_pwauth"`
	DisableRoot bool `yaml:"disable_root"`
}

type Network struct {
	Network struct {
		Version   int
		Ethernets struct {
			Ens192 struct {
				Match struct {
					Name string
				}
				Addresses   []string
				Gateway4    string
				Nameservers struct {
					Addresses []string
				}
			}
		}
	}
}

type Metadata struct {
	Network         string `json:"network"`
	NetworkEncoding string `json:"network.encoding"`
}

func GenerateBaseTemplate(sshKey string) *Template {
	template := Template{}

	template.PackageUpdate = true
	template.PackageUpgrade = true

	template.Users = []User{{
		SshAuthorizedKeys: []string{sshKey},
	}}

	template.Chpasswd.Expire = true

	template.Growpart.Mode = "auto"
	template.Growpart.Devices = []string{"/"}
	template.Growpart.IgnoreGrowrootDisabled = false
	template.SshPwauth = true
	template.DisableRoot = false

	return &template
}

func AddSpecificParameters(specifier string, template *Template, pass string, networkTemplate *Network) (*Template, *Metadata) {
	newTemplate := template

	newTemplate.Users[0].Name = "root"
	newTemplate.Chpasswd.List = []string{"root" + ":" + pass}

	if specifier == "ubuntu" {
		networkTemplate, _ := yaml.Marshal(networkTemplate)
		newTemplate.WriteFiles = []WriteFile{{Encoding: "base64", Content: base64.StdEncoding.EncodeToString(networkTemplate), Path: "/etc/netplan/50-cloud-init.yaml"}}
		newTemplate.Runcmd = []string{"echo \"PermitRootLogin yes\" >> /etc/ssh/sshd_config", "systemctl restart ssh", "netplan apply"}
	} else if specifier == "centos" {
		networkTemplate, _ := yaml.Marshal(networkTemplate.Network)
		metadata := Metadata{
			Network:         base64.StdEncoding.EncodeToString(networkTemplate),
			NetworkEncoding: "base64",
		}
		return newTemplate, &metadata
	}
	return newTemplate, &Metadata{}
}

func CreateNetworkTemplate(identifier string, ipToAssign string) *Network {
	template := Network{}
	template.Network.Version = 2
	ens192 := template.Network.Ethernets.Ens192
	ens192.Match.Name = "ens*"
	ens192.Addresses = []string{ipToAssign + "/24"}
	ens192.Gateway4 = os.Getenv(identifier + "_GATEWAY")
	ens192.Nameservers.Addresses = strings.Split(os.Getenv(identifier+"_NAMESERVERS"), ",")
	template.Network.Ethernets.Ens192 = ens192
	return &template
}
