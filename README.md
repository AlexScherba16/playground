# Hi, welcome to "playground", or little tool crafted in Golang 🚀

![](./docs/media/logo.png)
<!-- TOC -->
* [📖 General info](#-general-info)
* [💻 System requirements](#-system-requirements)
* [🏛️ Source code structure](#-source-code-structure)
* [⚙️ Build & Run](#-build--run)
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
* * [cmd/playground](cmd/playground) - entry point for the "playground" executable
* [docs/](docs) - documentation-related files for the project
* * [docs/license](docs/license) - project license agreement
* * [docs/media](docs/media) - project images
* * [docs/testdata](docs/testdata) - data samples used for demo and tests
* [internal/](internal) - internal packages that are not intended for external use
* * [internal/cli](internal/cli) - cli parser entity and tests
* * [internal/constants](internal/constants) - project constant variables
* * [internal/runners](internal/runners) - runners are entities that operate as goroutines in a data processing pipeline
* * * [internal/runners/common](internal/runners/common) - common runners interface
* * * [internal/runners/datasource](internal/runners/datasource) - data pipeline entry point, it provides records(raw data) to other runners
* * * * [internal/runners/datasource/csv](internal/runners/datasource/csv) - csv file runner implementation and tests
* [internal/types](internal/types) - structures and channels types for internal usage across the project
* * [internal/utils](internal/utils) - utility functions and helpers for internal usage across the project
* * * [internal/utils/cerror](internal/utils/cerror) - custom error handler, provides common error message template
* * * [internal/utils/parser](internal/utils/parser) - files data parser, converts file lines to records

## 🏗 Build & Run

## 📱 Contacts
``` 
email:      alexscherba16@gmail.com
telegram:   @Alex_Scherba
```
