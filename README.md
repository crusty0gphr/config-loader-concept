## Config loader concept architecture for GoLang services

This is a monorepository project that demonstrates a working concept of an architecture with services capable of hot-swapping included configurations.
___
To achieve the automated configuration update process, I have designed a system where the config-loader service acts as a centralized configuration distributor using NATS as the messaging system. Here’s a step-by-step outline of the architecture:

* **Config-Loader Service**:
  - Acts as the producer.
  - Reads configuration files.
  - Publishes configuration updates to NATS subjects.
* **Service Subscribers**:
  - Act as consumers.
  - Subscribe to specific NATS subjects to receive their respective configuration updates.
  - Implement hot-swap logic to apply the new configuration without restarting.
___
### Architecture Diagram
```
+-----------------+       +-----------------+       +-----------------+
|                 |       |                 |       |                 |
|  Service 1      |       |  Service 2      |       |  Service 3      |
|                 |       |                 |       |                 |
|   (Subscriber)  |       |   (Subscriber)  |       |   (Subscriber)  |
+-------+---------+       +-------+---------+       +-------+---------+
        |                         |                         |
        +-------------------------+-------------------------+
                                  |
                                  v
                          +-----------------+
                          |                 |
                          |  NATS Server    |
                          |                 |
                          +-------+---------+
                                  |
                                  v
                          +-----------------+
                          |                 |
                          | Config-Loader   |
                          |    (Producer)   |
                          |                 |
                          +-----------------+
```
___
### Basic implementation
* **Config-Loader Service**:
  - Reads YAML configuration files.
  - Publishes configuration updates to NATS subjects, e.g., config.service1, config.service2, etc.
* **Service Subscribers**:
  - Subscribe to their specific NATS subject to receive configuration updates.
  - Implement logic to apply the new configuration without restarting.
___
### Project directory scheme
```
root/
├── configloader
├── configstor/
│   ├── cfg.service1.yaml
│   ├── cfg.service2.yaml
│   └── cfg.service3.yaml
├── servicepool/
│   ├── service1/
│   │   └── config/
│   │       └── cfg.base.yaml
│   ├── service2/
│   │   └── config/
│   │       └── cfg.base.yaml
│   └── service3/
│       └── config/
│           └── cfg.base.yaml
├── docker-compose.yaml
└── Makefile
```

### Directory and File Descriptions
* **configloader/**: Contains the source code and related files for the config-loader service, responsible for distributing updated configuration files to other services.
* **configs/**: Stores the centralized configuration files that are managed and distributed by the config-loader service. Each YAML file corresponds to a different service:
  - cfg.service1.yaml: Configuration for service1
  - cfg.service2.yaml: Configuration for service2
  - cfg.service3.yaml: Configuration for service3
* **servicepool/**: Contains directories for each service with their initial base configurations:
  - **service1/**:
    - config/: Contains the base configuration file (cfg.base.yaml) for service1.
  - **service2/**:
    - config/: Contains the base configuration file (cfg.base.yaml) for service2.
  - **service3/**:
    - config/: Contains the base configuration file (cfg.base.yaml) for service3.