# Hi, welcome to "playground", or little tool crafted in Golang 🚀

![](./docs/media/logo.png)
<!-- TOC -->

* [📖 What's it all about?](#-whats-it-all-about)
* [💻 System requirements](#-system-requirements)
* [🏛️ Source code structure](#-source-code-structure)
* [🏗 Build & Run](#-build--run)
* [📱 Contacts](#-contacts)
<!-- TOC -->


## 📖 What's it all about?
``` 
Project goals and purposes should be here, but its confidential, sorry.
All that I can say, it something about data prediction, on N-th day.
Well, that's it. Cool, right?🤓 You can continue reading 👇
``` 
## 💻 System requirements
* go version go1.21.0

## 🏛️ Source code structure
* [cmd/](cmd) - application entry points directory
* * [playground](cmd/playground) - entry point for the "playground" executable
* [docs/](docs) - documentation-related files for the project
* * [license](docs/license) - project license agreement
* * [media](docs/media) - project images
* * [testdata](docs/testdata) - data samples used for demo and tests
* [internal/](internal) - internal packages that are not intended for external use
* * [cli](internal/cli) - cli parser entity and tests
* * [constants](internal/constants) - project constant variables
* * [runners/](internal/runners) - runners are entities that operate as goroutines in a data processing pipeline
* * * [aggregator/](internal/runners/aggregator) - data aggregator runners backed by a provided aggregation parameter
* * * * [aggregator_factory](internal/runners/aggregator/aggregator_factory) - aggregator runner creator and tests
* * * * [runner](internal/runners/aggregator/runner) - aggregator runner implementation and tests
* * * * [strategy/](internal/runners/aggregator/strategy) - data aggregation algorithms and tests
* * * * * [campaign](internal/runners/aggregator/strategy/campaign) - campaign data aggregation algorithm and tests
* * * * * [country](internal/runners/aggregator/strategy/country) - country data aggregation algorithm and tests
* * * [common](internal/runners/common) - common runners interface
* * * [datasource/](internal/runners/datasource) - data pipeline entry point, runners provide records(raw data) to other runners
* * * * [datasource_factory](internal/runners/datasource/datasource_factory) - datasource runner creator and tests
* * * * [runner/](internal/runners/datasource/runner) - datasource runners implementation
* * * * * [csv](internal/runners/datasource/runner/csv) - csv file runner implementation and tests
* * * * * [json](internal/runners/datasource/runner/json) - json file runner implementation and tests
* * * [postprocessor/](internal/runners/postprocessor) - final part of data pipeline, prepares predicted data to console output
* * * * [postprocessor_factory](internal/runners/postprocessor/postprocessor_factory) - postprocessor runner creator and tests
* * * * [runner](internal/runners/postprocessor/postprocessor_factory) - postprocessor runner creator and tests
* * * * [strategy/](internal/runners/postprocessor/strategy) - postprocessor algorithms and tests
* * * * * [campaign](internal/runners/postprocessor/strategy/campaign) - campaign data postprocessor algorithm and tests
* * * * * [country](internal/runners/postprocessor/strategy/country) - country data postprocessor algorithm and tests
* * * [predictor/](internal/runners/predictor) - data predictor runners backed by a provided model parameter
* * * * [predictor_factory](internal/runners/predictor/predictor_factory) - predictor runner creator and tests
* * * * [runner](internal/runners/predictor/runner) - predictor runner implementation and tests
* * * * [strategy/](internal/runners/predictor/strategy) - predictor data algorithms
* * * * * [linext](internal/runners/predictor/strategy/linext) - linear extrapolation data predictor and tests
* * * * * [average](internal/runners/predictor/strategy/average) - average data predictor and tests
* [types](internal/types) - structures and channels types for internal usage across the project
* * [utils/](internal/utils) - utility functions and helpers for internal usage across the project
* * * [cerror](internal/utils/cerror) - custom error handler, provides common error message template
* * * [parser](internal/utils/parser) - files data parser, converts file lines to records
* * * [predictor](internal/utils/predictor) - predictor algorithms util functions, math stuff

## 🏗 Setup & Run
``` 
git clone https://github.com/AlexScherba16/playground
cd playground
go run cmd/playground/main.go -source docs/testdata/test_data.csv -model linext -aggregate country

Enjoy 😉
```

## 📱 Contacts
``` 
email:      alexscherba16@gmail.com
telegram:   @Alex_Scherba
```
