# Hobbit Registry
Automation for downloading and pushing private registry images

## Settings
* Create a yaml file with the configs <br>

Example:
```yaml
configs:
  registry:
    scheme: "http"
    url: "private repo" # private registry
    port: 5000
    #username: "username"
    #password: "password"
  images: # required images array
    - "redis:6.2.4"
```

## How to use it?
```bash
$ git clone && cd hobbit-registry
$ go install
$ hobbit-registry -c <PATH to config_file.yaml>
```