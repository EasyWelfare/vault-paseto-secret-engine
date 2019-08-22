# vault-paseto-secret-engine

Valt secret engine plugin for Paseto


## Enable Plugin

Due to the fact that is a custom dynamic secret engine, is not enabled by default, so you have to enable it manually with the following command:

```bash
vault secrets enable -path=<path> paseto
```

Where `<path>` is the path where will be activated the dynamic secret engine plugin

## Plugin Configuration

This plugin needs to be configured firstly, with two simple parameters: 

1. ttl: the expiration time of each token that will be released at every token invocation
2. footer: a discriminant parameter that is appended to the token

This kind of setup is done through the following command : 


```bash
vault write <path/paseto/config> footer="whatever" ttl=120
```

to check the just written configuration you can hit:

```bash
vault read <path/paseto/config>
```

## Exposed Entrypoints

Once the plugin is configured, accordingly to the ACL set (see ACL section below), you can access to :

1. <path/paseto/token> : if you are allowed to access to this path, it will return a paseto compliant token with the expiration time set in the `ttl` configuration previously set

2. <path/paseto/public> : if you are allowed to access to this path, it will return the public key

3. <path/paseto/config> : if you are allowed to access to this path, it will give you  the configuration previously set

the `<path>` is intended to match with the configuration path previously set

## ACL policies

As requested we makes 2 roles:

1. _owner_ : (aka the server) will have access to the configuration and public
2. _clients_: will have access to the token entrypoints