# Architectural Requirements Document

<p align="center">
  <img src="assets/spark_logo_dark_background.png" alt="Company Logo" width="300" height=100%/>
</p>

**Version:** 3.0.0.0  
**Prepared By:** Spark Industries  
**Prepared For:** Southern Cross Solutions  
**Document Type:** Architectural Requirements Document  
**Demo:** Demo 3

## Table of Contents

- [Architectural Quality Requirements](#architectural-quality-requirements)
- [Architectural Strategies](#architectural-strategies)
- [Architectural Design and Pattern](#architectural-design-and-pattern)
- [Architectural Constraints](#architectural-constraints)

---


## Architectural Quality Requirements

The following quality requirements are prioritized from highest to lowest priority for the Smart File Manager system:

## Architectural Quality Requirements for Smart File Manager (SFM)

### 1. **Reliability**

- **Specification:** The system must maintain file integrity and recover gracefully from failures (e.g., crashes or power loss).
- **Testable Criteria:** No data corruption or loss during operations; system auto-recovers to last stable state after unexpected shutdowns.
- **Rationale:** Users must trust the system to manage personal or critical files safely and without risk of accidental loss.

### 2. **Performance**

- **Specification:** File operations, such as sorting, tagging, and smart search, must execute within acceptable time limits on standard desktop hardware.
- **Testable Criteria:** Smart search results must appear in under 2 seconds; bulk operations (e.g., tagging or reclassifying 500 files) must complete within 5 seconds.
- **Rationale:** Responsiveness is key to a smooth user experience in a personal desktop application. Delays can disrupt workflow and reduce trust in automation.

### 3. **Scalability**

- **Specification:** The system must handle increasing numbers of files, metadata, and tags without degradation in performance.
- **Testable Criteria:** Maintain acceptable performance (search in <2 seconds, classification in <5 seconds) when managing up to 1 million files or 1TB of data.
- **Rationale:** Personal file collections can grow substantially over time, so the system must remain efficient and responsive even at high volume.

### 4. **Usability**

- **Specification:** The user interface must be intuitive, accessible, and require minimal onboarding or training.
- **Testable Criteria:** First-time users can complete core tasks (e.g., find a file, create a smart manager, apply a tag) within 5 minutes, with no external documentation.
- **Rationale:** The system targets general users with varying technical skill levels; high usability promotes adoption and continued use.

### 5. **Modifiability**

- **Specification:** Users must be able to define or update rules, filters, and semantic tags without developer intervention.
- **Testable Criteria:** 100% of common modifications can be made through a graphical user interface without restarting the application.
- **Rationale:** Enables users to customize the system to suit their unique workflows and evolving organizational habits.

---


## Architectural Strategies

<table>
  <thead>
    <tr>
      <th>Quality Attribute</th>
      <th>Architectural Strategy</th>
      <th>Implementation Details</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Reliability</td>
      <td>Transactional File Operations</td>
      <td>All operations on the user files should either work completely for the entire operation or not modify the user's files at all.</td>
    </tr>
    <tr>
      <td>Performance</td>
      <td>Concurrency</td>
      <td>Leverage our master slave pattern to break slow requests into smaller operations which may be run concurrently. </td>
    </tr>
    <tr>
      <td>Performance</td>
      <td>I/O overhead</td>
      <td>Minimize the number of times files need to be opened to reduce I/O overhead. </td>
    <tr>
      <td>Scalability</td>
      <td>Horizontal Scaling</td>
      <td>Create more instances of Python slaves to increase throughput for high demand. </td>
    </tr>
    <tr>
      <td>Usability</td>
      <td>User-centered design</td>
      <td>Conduct usability tests, provide intuitive navigation, and minimize complexity. Provide accessible interfaces appropriate for users of varying skill level</td>
    </tr>
    <tr>
      <td>Modifiability</td>
      <td>Modular, loosely coupled architecture</td>
      <td>Design the system such that all operations accept various parameters while still having default ones. Allow users to customize these parameters from the user-interface.</td>
    </tr>
  </tbody>
</table>


---

## Architectural Design and Pattern

SFM makes use of various architectural patterns selected to support our quality requirements. Our choices and justifications for them follows.

### Monolithic Architecture 

Monolithic architecture combines the various components of an application into a single inseparable unit. Our project is deployed as an single executable application running exclusively on the user's personal machine, with no external dependancies. This deployment model justifies our choice for using this architecture.

#### How it supports our quality requirements
* **Reliability:**
   - Monolithic systems have no dependancies on external modules which could fail, therby reducing overall system uptime.
   - Monolithic systems are not continuously deployed, instead only being deployed once a stable and fully tested version is ready, improving reliability.
* **Performance:**
   - Due to reduced communication overhead between modules within the monolithic system performance may be faster as opposed to network based bottlenecks introduced by other architectures.
* **Scalability:**
   - While monolithic systems do not scale horizontally, monolithic systems may scale vertically with improved hardware. Please also see how we leverage the master-slave pattern to offset the limited scalability. 
* **Usability:**
   - By deploying our monolithic system on different operating systems we improve usability since the application runs in an environment familiar to the end user. 
* **Modifiability:**
   - Not applicable 

### Master-Slave Architectural Pattern

Master-Slave architecture leverages a central node (master) to delegate tasks to multiple worker nodes (slaves). We mainly leverage this pattern to improve the performance of our computational demanding tasks.

#### How it supports our quality requirements
* **Reliability:**
   - Reliance on a master node generally decreases the reliability of the master-slave architecture pattern. We compensate for this by ensuring that when our master fails it does so gracefully without going into an unrecoverable state. Use of the monolithic architecture further supports this since the domain in which our system runs is more predictable (owing to the deployment).
* **Performance:**
   - The primary reason for using Master-slave is performance. Our backend is paritioned into 2 sub-systems reponsible for filesystem and clustering features respectively. By leveraging threadpools in our clustering subsystem (slaves) we increase performance drastically.
   - Furthermore our slaves are stateless and do not interact with the same resources, avoiding traditional concurrency related problems like deadlock.
* **Scalability:**
   - The Master-slave patterns facilitates horizontal scaling by creating more slave instances. It should be noted that there is a upper limit on our ability to scale horizontally dictated by the system's number of logical processors.
* **Usability:**
   - Not applicable
* **Modifiability:**
   - Not applicable 



### System Components

#### Website Component

- User-friendly interface for application download
- Account registration and management functionality

#### Application Component

- File explorer functionality
- File/folder sorting by tags and metadata
- Structure graph visualization
- Advanced search capabilities
- File and folder management operations

#### Sub-Architectural Patterns

- **Layered Architecture (N-tier):** Separates concerns into Presentation, Application, and Data layers

![Website Architecture](assets/websiteArchitecure.png)
![Application Architecture](assets/applicationArchitecture.png)

---

## Architectural Constraints

The following constraints affect the architectural design of the Smart File Manager:

### Data Constraints

- File system access - Which components can access which directories or file types on the system
- Data persistance - Where can files be stored and naming conventions
- File format - restrictions on file types and accessing metadata.

### Technology Constraints

- Operating system compatibility - Must run on Windows, macOS, Linux platforms, thus need architecture that supports cross-platform.

### Performance Constraints

- runtime bottleneck - Slower execution affects real-time file operations.
- overhead in memory - higher memory usage for large file sorting operations.

### Deployment Constraints

- The system must be deployable as a desktop app
- Must run on Windows, Linux, and macOS with consistent functionality.

---

