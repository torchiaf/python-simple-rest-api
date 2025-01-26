- Modules, Devices, Circuit -> should be created by Crd POST, not helm install
  - We could also use Fleet
- Rename adapter -> rpc_client
- Implement OpenAPI definitions
- Golang process on master node to call the rabbitMQ endpoints, fetch data and update CRDS
- Create Rancher extensions to handle sensors CRDs
- Create a manifest.yaml file for each supported sensors (svg, description, link to homepage)
- Sensor's CRDS should be updated ONLY by controller pod
- Rancher UI should display
  - sensor CRD
  - raspberry CRD
- Define a builder kit for devices
  - The device <-> rpc-server interface should use OpenAPI definition to build a skeleton and call the executable built in devices docker images.
  - The device's API should be defined in the settings file by dev
- Create Circuits to connect devices
- Define job python code using code-server in rancher extension 
- Inject python script in Job's ConfigMap

- Go dependecy injection do define lmbda to apply to Circuits 
  https://medium.com/avenue-tech/dependency-injection-in-go-35293ef7b6
  Google Wire https://github.com/google/wire?tab=readme-ov-file

- Go Plugins to define Circuit workloads
- Circuit crd should have play, stop, pause fields
- Circuit crd should have a base64 to store code-source