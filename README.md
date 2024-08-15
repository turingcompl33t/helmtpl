## `helmtpl`

`helmtpl` (_Helm Template_) is a tool for applying templating logic to Helm _values.yaml_ configuration files.

### Usage

The input to `helmtpl` is a _template file_. The template file utilizes YAML syntax the same way a Helm `values.yaml` file does. This template files declares all of the same configuration data that you want to include in your `values.yaml` file. In addition, it also includes a special key, `vars`, which can be used to inject variables into the configuration files. `helmtpl` then resolves these variable references during processing to create your final `values.yaml` file which will be used as input to Helm.

Suppose we have the following `values.tpl.yaml` file:

```yaml
# values.tpl.yaml
vars:
  suba:
    name: suba

suba:
  name: "{{ .vars.suba.name }}"

subb:
  subaname: "{{ .vars.suba.name }}"
```

Then we can invoke `helmtpl`:

```bash
helmtpl -input values.tpl.yaml
```

To produce the output file `values.yaml`:

```yaml
# values.yaml
suba:
    name: suba
subb:
    subaname: suba
```
