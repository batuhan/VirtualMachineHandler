package helpers

import (
	"encoding/base64"
	"github.com/sethvargo/go-password/password"
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
	WriteFiles []WriteFile `yaml:"write_files"`
	Runcmd     []string
}

type Network struct {
	Network struct {
		Version   int
		Ethernets struct {
			Ens192 struct {
				Addresses   []string
				Gateway4    string
				Nameservers struct {
					Addresses []string
				}
			}
		}
	}
}

func GenerateBaseTemplate(sshKey string) *Template {
	template := Template{}

	template.PackageUpdate = true
	template.PackageUpgrade = true

	template.Users = []User{{
		SshAuthorizedKeys: []string{sshKey},
	}}

	template.Chpasswd.Expire = false

	template.Growpart.Mode = "auto"
	template.Growpart.Devices = []string{"/"}
	template.Growpart.IgnoreGrowrootDisabled = false

	return &template
}

func AddUbuntuSpecificParameters(template *Template, networkTemplate []byte) (*Template, error) {
	newTemplate := template

	newTemplate.Users[0].Name = "ubuntu"
	pass, err := password.Generate(12, 2, 2, false, false)
	if err != nil {
		return nil, err
	}
	newTemplate.Chpasswd.List = []string{"ubuntu:" + pass}

	newTemplate.WriteFiles = []WriteFile{{Encoding: "base64", Content: base64.StdEncoding.EncodeToString(networkTemplate), Path: "/etc/netplan/50-cloud-init.yaml"}}
	newTemplate.Runcmd = []string{"netplan apply"}
	return newTemplate, nil
}

func CreateNetworkTemplate(identifier string, ipToAssign string) *Network {
	template := Network{}
	template.Network.Version = 2
	ens192 := template.Network.Ethernets.Ens192
	ens192.Addresses = []string{ipToAssign + "/24"}
	ens192.Gateway4 = os.Getenv(identifier + "_GATEWAY")
	ens192.Nameservers.Addresses = strings.Split(os.Getenv(identifier+"_NAMESERVERS"), ",")
	template.Network.Ethernets.Ens192 = ens192
	return &template
}
