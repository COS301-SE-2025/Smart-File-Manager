<p align="center">
  <img src="assets/spark_logo_dark_background.png" alt="Company Logo" width="300" height=100%/>
</p>

# Coding Standards Document for SparkIndustries

**Version:** 1.0.0.0  
**Prepared By:** Spark Industries  
**Prepared For:** Personal Development team use  

## Content
* [Introduction](#introduction)
* [Naming Conventions](#naming-conventions)
* [Glossary](#glossary)

## Introduction 
The following document serves as the official coding standards document for SparkIndustries. It outlines all best practices and conventions our team will use to ensure the maintenance of a well formed, readable and uniform codebase. Going forward our team will aim to use this document as guidelines for personal development and during our PR review process. 
   
It is important to take note that our codebase consists of three primary parts each developed in its own language. These are Python, Go and Flutter (Dart). Conventions for these languages may differ. In cases where this observation holds we decide on choosing the conventions accepted by the individual languages rather than trying to force them across multiple programming domains. In any scenarios where this does not hold we provide ample justification for our choices.


## Naming Conventions
<table>
  <thead>
    <tr>
      <th>Concept</th>
      <th>Python</th>
      <th>Go</th>
      <th>Flutter (Dart)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><strong>Variables</strong></td>
      <td>snake_case</td>
      <td>camelCase</td>
      <td>camelCase</td>
    </tr>
    <tr>
      <td><strong>Functions / Methods</strong></td>
      <td>snake_case (Use type hinting)</td>
      <td>camelCase (unexported), PascalCase (exported)</td>
      <td>camelCase</td>
    </tr>
    <tr>
      <td><strong>Constants</strong></td>
      <td>ALL_CAPS</td>
      <td>PascalCase (exported), camelCase (unexported)</td>
      <td>camelCase </td>
    </tr>
    <tr>
      <td><strong>Classes / Types</strong></td>
      <td>PascalCase</td>
      <td>PascalCase</td>
      <td>PascalCase</td>
    </tr>
    <tr>
      <td><strong>File Names</strong></td>
      <td>snake_case.py</td>
      <td>camelCase.go</td>
      <td>snake_case.dart</td>
    </tr>
    <tr>
      <td><strong>Packages / Modules</strong></td>
      <td>snake_case</td>
      <td>lowercase (no underscores)</td>
      <td>snake_case</td>
    </tr>
    <!-- Add more rows as needed -->
  </tbody>
</table>

In general we also propose the following conventions:
- Use descriptive variable, class and method names.
- For python: Make explicit use of type hinting.
- Do not use ambigous abbreviations. 
- Avoid variable hiding in nested scopes.

## Indentation
Python, Go and Flutter should all use tabs. This holds true for all use of indentation.

## File Layout
#### Python

```
.
└── Python/
    ├── src/
    │   ├── file1.py
    │   ├── file2.py
    │   └── file3.py
    └── testing  /
        ├── test_file1.py
        ├── test_file2.py
        └── test_files/
            ├── testimg.jpeg
            └── testdoc.doxc
```

### Go
```
.
└── golang/
    ├── client/
    │   └── all the generated protos
    ├── filesystem/
    │   ├── goFile1.go
    │   └── goFile2.go
    └── gprc/
        ├── client/
        │   └── grpcClient.go
        └── server/
            └── grpcServer.go
```

### Flutter
```
.
└── app/
    ├── images
    ├── lib/
    │   ├── api.dart
    │   ├── constants.dart
    │   └── main.dart
    ├── linux/
    │   ├── flutter
    │   └── runner/
    │       ├── *.cpp
    │       └── *.h
    ├── macOS/
    │   ├── flutter
    │   └── runner/
    │       └── *.swift
    ├── test/
    │   └── widget_test.dart
    └── windows/
        ├── flutter
        └── runner/
            ├── *.cpp
            └── *.h
```

## Exception Handling
We do not enforce a uniform way to handle exceptions throughout our system due to the varying nature of exception handling in the different project technologies. Developers should follow the best practices associated with their language.  
In general we propose the following guidelines:
- Use specific exception types (i.e. IndexOutOfBounds vs RuntimeError)
- Handle Exceptions at as low a level as reasonably possible. Do not continously propogate exceptions up the call-stack.
- Do not "swallow" exceptions. 
- Log exceptions with clear and descriptive messages.
- Avoid using exceptions for flow control
- Be mindful of logging sensitive information in error messages. 

## Commenting
We propose the following guidelines for commenting:
- Comment complex logic instead of blindly commenting what can be inferred from method names
- Make use of TODO comments where applicable
- Where multiple approaches could be taken, use comments to explain rationale for choosing one.
- Do not commit large commented out blocks. Make use of version control for its intended use. 
- Explain **why** not **what**

## General
The following general guidelines must be kept in mind:
- Do not implement functionality already present in a used library.


## Glossary
In an effort to keep this document concise we do not belabour what we mean by well known naming conventions. We briefly explain them here for easy reference.

* camelCase: Entails combining multiple words into a single word with the first letter of each word (except the very first word) capitalized. E.g. myName, fooBar.
* snake_case: Entails combining multiple words into a single word separated by underscores. E.g. my_name, foo_bar.
* PascalCase: Entails combining multiple words into a single word with the first letter of each word capitalized.
