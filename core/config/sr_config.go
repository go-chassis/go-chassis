package config

import (
	"fmt"
	"github.com/go-chassis/cari/rbac"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
	"github.com/go-chassis/openlog"
)

// GetRegistratorType returns the Type of service registry
func GetRegistratorType() string {
	return GlobalDefinition.ServiceComb.Registry.Type
}

// GetRegistratorAddress returns the Address of service registry
func GetRegistratorAddress() string {
	return GlobalDefinition.ServiceComb.Registry.Address
}

// GetRegistratorScope returns the Scope of service registry
func GetRegistratorScope() string {
	return GlobalDefinition.ServiceComb.Registry.Scope
}

// GetRegistratorAutoRegister returns the AutoRegister of service registry
func GetRegistratorAutoRegister() string {
	return GlobalDefinition.ServiceComb.Registry.AutoRegister
}

// GetRegistratorAPIVersion returns the APIVersion of service registry
func GetRegistratorAPIVersion() string {
	return GlobalDefinition.ServiceComb.Registry.APIVersion.Version
}

// GetRegistratorDisable returns the Disable of service registry
func GetRegistratorDisable() bool {
	return archaius.GetBool("servicecomb.registry.disabled", false)
}

// GetRegistratorRbacAccount returns the RbacAccout info of service registry
func GetRegistratorRbacAccount() *rbac.AuthUser {
	username := GlobalDefinition.ServiceComb.Credentials.Account.Username
	password := GlobalDefinition.ServiceComb.Credentials.Account.Password
	cipherName := GlobalDefinition.ServiceComb.Credentials.AccountCustomCipher
	cipherPlugin, err := cipher.NewCipher(cipherName)
	if err != nil {
		openlog.Error(fmt.Sprintf("get cipher plugin [%s] failed, %v", cipherName, err))
		return nil
	} else if cipherPlugin == nil {
		openlog.Error(fmt.Sprintf("invalid cipher plugin"))
		return nil
	}

	pwd, err := cipherPlugin.Decrypt(password)
	if err != nil {
		openlog.Error(fmt.Sprintf("Decrypt account password failed %v", err))
		return nil
	}

	auth := rbac.AuthUser{
		Username: username,
		Password: pwd,
	}
	return &auth
}
