# Figs

## Introduction
Figs is short for `Fi`le `G`eneration `S`ervice. It is a tool that allows you to generate files from templates.

Figs converts data inputs into document output based on predefined configuration and templates. 

## Components
The service consists of, at least, two components:
- Core Library - this a pure go library that expose simple, high-level APIs for consumption by users.
- HTTP Wrapper - A RESTful service that makes use of the core library interface to generate documents based on requests. 
- [Optional] CLI Wrapper - A CLI wrapper around the core library  

## Core Library
### Configuration
- Data Mapping Schema
- Data Validation Schema
- Document Storage Driver/Settings

### Data Input
- Input data to the core service should contain three key information:
- Type - This identifier determines which validation and data mapping schema would be use against the payload
- Payload - The actual data to be converted
- Document Type - either CSV or PDF. Both could also be generated at the same time

### Document Output
- The service is capable of generating CSV and PDF documents. PDF files generated is base on predefined templates. Templates can in the default go template format or a cross-platform templating format like mustache. Defining template using mustached has the advantage of portability for other purpose.

- Every template would also have an identifier. This identifier links a request, validation schema, data mapping function and template. 

- Generated files would be sent to a configurable external storage service like s3, ftp etc

## Engine
The data to document engine consists of three main parts:
- Validator
- Mapper
- Converter

### Validator
- Takes in the request payload, and the type identifier
- Uses the type identifier to select a validation schema
- Validates the payload with the selected schema
- If validation passes, moves to the next step. Else throws an error

### Mapper
- Takes in the request payload, and the type identifier
- Uses the type identifier to select a mapping function
- Maps the payload to the correct/standard document data format
- If mapping succeeds, moves to the next step. Else throws an error

### Converter
- Takes in the request payload, and the document type
- Uses the document type to select the right converter function, CSV or PDF
- Generates file, store the file and return a unique URL for the file

## Usage
- Use CLI
```bash
$ figs generate \
    --data=[a file path] \
    --template=[a pre-configured template name] \
    --format=[pdf,excel] \
    --storage=[local,s3] 
// generated files would be stored in the configured storage location and available via http request for download
```

- Use HTTP
```bash
// First start the service
$ figs server start --port=8080

// Then make a request
$ curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "data": {...}, 
        "template":"[a pre-configured template name]", 
        "format":"[pdf,excel]", 
        "storage":"[local,s3]"
        }' \
    http://localhost:8080/generate

// generated files url would be returned
```

## Data Schema
    - segments with predefined data types
    - - header
    - - signature
    - - letter 
    - - body
    - - table
    - - etc, extendable 
    - select predefined template (mustache based)
    - just data
    - output types pdf, excel
    - excel would only use table

## Notes
Interface should be well defined so the can be easily extensible. For instance, it should be easy to plug-in a different storage engine or document convert or data mapper

