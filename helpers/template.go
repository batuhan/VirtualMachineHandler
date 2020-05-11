package helpers

import (
	"encoding/base64"
	"gopkg.in/yaml.v2"
	"strings"
)

type User struct {
	Name              string
	SshAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
}

type WriteFile struct {
	Encoding    string
	Content     string
	Path        string
	Permissions string
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

type Ethernet struct {
	Match struct {
		Name string
	}
	Addresses   []string
	Gateway4    string
	Nameservers struct {
		Addresses []string
	}
}

type Network struct {
	Network struct {
		Version   int
		Ethernets struct {
			Eth0   Ethernet
			Ens192 Ethernet
		}
	}
}

type Metadata struct {
	Network         string `json:"network"`
	NetworkEncoding string `json:"network.encoding"`
}

func GenerateBaseTemplate(sshKeys []string, onFirstBoot []string) *Template {
	template := Template{}

	template.PackageUpdate = false
	template.PackageUpgrade = false

	template.Users = []User{{
		SshAuthorizedKeys: sshKeys,
	}}

	template.Runcmd = onFirstBoot

	template.Chpasswd.Expire = true

	template.Growpart.Mode = "auto"
	template.Growpart.Devices = []string{"/dev/sda2", "/dev/sda5"}
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
		newTemplate.WriteFiles = []WriteFile{
			{Encoding: "base64", Content: base64.StdEncoding.EncodeToString(networkTemplate), Path: "/etc/netplan/01-netcfg.yaml"},
			{Encoding: "base64", Content: base64.StdEncoding.EncodeToString([]byte("#!/bin/sh\npvresize /dev/sda5\nlvresize -l +100%FREE /dev/mapper/vg-root\nresize2fs /dev/mapper/vg-root\n")), Path: "/var/lib/cloud/scripts/per-boot/diskresize.sh", Permissions: "755"}}
		newTemplate.Runcmd = append(newTemplate.Runcmd, "echo \"PermitRootLogin yes\" >> /etc/ssh/sshd_config", "systemctl restart ssh", "netplan apply")
	} else if strings.Contains(specifier, "centos") || specifier == "debian" {
		var vgName string
		if specifier == "centos-7" {
			vgName = "centos"
		} else if specifier == "centos-8" {
			vgName = "cl"
		} else if specifier == "debian" {
			newTemplate.Runcmd = append(newTemplate.Runcmd, "perl -pe 'BEGIN{undef $/;} s|iface ens.*||gs' /etc/network/interfaces > /etc/network/interfaces.new", "mv /etc/network/interfaces.new /etc/network/interfaces", "ip addr flush dev ens192", "service networking restart")
			vgName = "localhost--vg"
		}

		newTemplate.WriteFiles = []WriteFile{
			{Encoding: "base64", Content: base64.StdEncoding.EncodeToString([]byte("#!/bin/sh\npvresize /dev/sda2\npvresize /dev/sda5\nlvresize -l +100%FREE --resizefs /dev/mapper/" + vgName + "-root\n")), Path: "/var/lib/cloud/scripts/per-boot/diskresize.sh", Permissions: "755"},
		}
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
	eth0 := template.Network.Ethernets.Eth0
	eth0.Match.Name = "ens*"
	eth0.Addresses = []string{ipToAssign + "/24"}
	eth0.Gateway4 = Config.Locations[identifier].Gateway
	eth0.Nameservers.Addresses = Config.Locations[identifier].Nameservers
	template.Network.Ethernets.Eth0 = eth0
	template.Network.Ethernets.Ens192 = eth0
	return &template
}
