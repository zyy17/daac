# 🚧 dac

`dac` is short for **Dashboards as Code**. It helps the developers to create the Grafana(for now) dashboards by a simple YAML file.

It's just a very simple wip project.

## How it works

The main idea is to **separate** the metric queries from the Grafana dashboard JSON model and define the more simple dashboard model in YAML file.
The `dac` will render the final Grafana dashboard JSON model by combining the two files.

The `dashboard.yaml` may look like this:

```yaml
name: my-dashboard
version: v1alpha1
title: "My Dashboard"
style:
   # TODO(zyy17): The reference also can be a URL.
   reference: styles/my-dashboard-style.json
groups:
   - name: http-requests
     title: "HTTP Requests"
     isRow: false
     panels:
        - name: http-requests
          title: "HTTP Requests"
          description: "HTTP requests per second"
          queries:
             - expr: "rate(myapp_http_requests_total[$__rate_interval])"
```

## How to use

1. Build the project

   ```console
   make
   ```
   
   Then you will get the `dac` binary in the `bin` directory.

2. Run the example

   ```console
   ./bin/dac -f examples/dashboard.yaml -o bin/dashboard.json
   ```

3. Use Grafana to import the `dashboard.json` file.

## Next

There's a lot of work to do:

- [ ] The more complicated render engine
- [ ] The layout system(maybe we don't always need to modify the panels manually)
- [ ] Define a general and flexible spec of the dashboard model
- [ ] Integrate with the Grafana API
- [ ] Integrate with the K8s encosystem
...
