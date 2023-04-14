package account

import "gonexa/config"

type Account struct {
	Mnenmonic     string
	PrivateKey    string
	PublicKey     string
	Address       string
	ChangeAddress string
}

func New(mnenmonic string, privatekey string, publickkey string, address string, changeAddress string) Account {
	return Account{mnenmonic, privatekey, publickkey, address, changeAddress}
}

func NewFromConfig(configAc config.MainConfig) Account {
	return Account{configAc.Mnenmonic, configAc.PrivateKey, configAc.PublicKey, configAc.Address, configAc.ChangeAddress}
}

func (account Account) ToString() string {

	return "Account{\n" +
		"Mnenmonic: " + account.Mnenmonic + "\n" +
		"PrivateKey: " + account.PrivateKey + "\n" +
		"PublicKey: " + account.PublicKey + "\n" +
		"Address: " + account.Address + "\n" +
		"ChangeAddress: " + account.ChangeAddress + "\n" +
		"}\n"
}

// 会自动调用
func init() {
	config.ConfigInit()
}

func GetMainAccount() Account {
	main := config.AllConfig()[0]
	return NewFromConfig(main)
}

func GetAccount(index uint8) Account {
	// config.ConfigInit()
	main := config.AllConfig()[index]
	return NewFromConfig(main)
}
