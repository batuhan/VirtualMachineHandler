package helpers

import "github.com/sethvargo/go-password/password"

type User struct {
	Name              string
	SshAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
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

	return &template
}

func AddUbuntuSpecificParameters(template *Template) (*Template, error) {
	newTemplate := template

	newTemplate.Users[0].Name = "ubuntu"
	pass, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return nil, err
	}
	newTemplate.Chpasswd.List = []string{"ubuntu:" + pass}
	return newTemplate, nil
}
