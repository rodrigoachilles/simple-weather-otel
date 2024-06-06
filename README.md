Translations:

* [French](README_fr.md)
* [Portuguese (Brazil)](README_pt_br.md)

# 🔍 Monitoring Project with OpenTelemetry and Zipkin (simple-weather-otel)

![Project Logo](assets/open-telemetry-zipkin.jpeg)

Welcome to the monitoring project with OpenTelemetry and Zipkin! This project demonstrates the use of OpenTelemetry for distributed tracing and monitoring of microservices in Go.

## 📑&nbsp;Table of Contents

- [📖 Introduction](#introduction)
- [🛠 Prerequisites](#prerequisites)
- [⚙️ Installation](#installation)
- [🚀 Usage](#usage)
- [🔎 Monitoring Examples](#monitoring-examples)
- [🤝 Contribution](#contribution)
- [📜 License](#license)

## 📖&nbsp;Introduction

OpenTelemetry is a collection of tools, APIs, and SDKs that can be used to instrument, generate, collect, and export telemetry data (such as metrics, logs, and traces) to help understand software behavior. Zipkin is used to collect and visualize this data.

## 🛠&nbsp;Prerequisites

Make sure you have the following items installed before continuing:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

Change the `docker-compose.yaml` file and add the `WeatherAPI` API key to query the desired temperatures (KEY_WEATHER_API):

- [WeatherAPI](https://www.weatherapi.com/)

## ⚙️&nbsp;Installation

1. Clone this repository:

    ```sh
    git clone git@github.com:rodrigoachilles/simple-weather-otel.git
    cd simple-weather-otel
    ```

2. Run Docker Compose:

    ```sh
    docker-compose up -d
    ```

3. Access Zipkin at:

   [http://localhost:9411](http://localhost:9411)

## 🚀&nbsp;Usage

After starting Docker Compose, you can access the Zipkin interface to monitor your service spans. To execute the services, use the `.http` file in the `api` folder of `weather-service`.

### 🔧&nbsp;Running Services

1. Navigate to the `api` folder in the `weather-service` directory:

    ```sh
    cd weather-service/api
    ```

2. Execute the `.http` file using your preferred tool (e.g., VSCode REST Client, Postman):

    ```sh
    # Example for VSCode REST Client
    weather.http
    ```

Here is an example of how a span can be visualized in Zipkin:

![Zipkin Span Example](assets/cep-service-spans.png)

### 💻&nbsp;Example Code

Here are some examples of how you can instrument your Go code to send data to Zipkin using OpenTelemetry:

#### Go Example

Install the necessary dependencies:

```sh
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/zipkin
go get go.opentelemetry.io/otel/sdk/trace
```

Instrument your application:

```go
package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
)

func main() {
	// Create Zipkin exporter
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
	)
	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			"service.name", "weather-service",
		)),
	)
	otel.SetTracerProvider(tp)

	// Your application code here
}
```

## 🔎&nbsp;Monitoring Examples

Below are examples of how spans from the `weather-service` and `cep-service` can be visualized in Zipkin:

![Weather Service Span](assets/weather-service-spans.png)
![CEP Service Span](assets/cep-service-spans.png)

## 🤝&nbsp;Contribution

Feel free to open issues or submit pull requests for improvements and bug fixes.

## 📜&nbsp;License

This project is licensed under the MIT License.
