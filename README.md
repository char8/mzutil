# mzutil

**work in progress**

A small utility for linux that does OAuth2 client flow, credential storage
and implements a simple Monzo API client. It's intended to be used in
conjunction with other linux tools, current targets:

* Integrate with [i3blocks](https://github.com/vivien/i3blocks) to show account balance in the status bar
* Integrate with [rofi](https://github.com/DaveDavenport/rofi) to show recent transactions on a hotkey

## TODO:

- [x] store secrets (OAuth token, secrets) on login keychain
- [x] `mzutil setup` - prompt for OAuth2 config
- [x] `mzutil login` - oauth2 login flow by opening browser and bringing up temp server for callback
- [x] `mzutil accounts` - list accounts
- [x] `mzutil balance` - print account balance
- [ ] `mzutil token` - print OAuth2 token expiry
- [ ] `mzutil tx` - list recent transactions
- [ ] Add scripts for rofi/i3blocks

## Uses:

- [skratchdot/open-golang](https://github.com/skratchdot/open-golang)
- [zalando/go-keyring](https://github.com/zalando/go-keyring)
- [golang.org/x/oauth2](https://github.com/golang/oauth2)
- [spf13/cobra](https://github.com/spf13/cobra)
