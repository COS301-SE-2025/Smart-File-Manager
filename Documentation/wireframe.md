# Smart File Manager - Wireframe Documentation

## Overview

This document provides a wireframe description of the Smart File Manager (SFM) application's user interface, including all main screens and their components.

## Global Navigation

![Sidebar Navigation](assets/wireframe_main_navigation.png)

- **Main Navigation Menu:**
  - Dashboard - Home screen with metrics and quick access
  - Smart Managers - Configuration page Smart Managers
  - Advanced Search - Enhanced file search for entire system
  - Settings - Application settings and options
- **Smart Managers List:**
  - Lists all current Smart Managers
  - Option to add new Smart Manager

## 1. Dashboard (Image 1)

![Dashboard View](images/dashboard.png)

### Purpose

Main overview screen displaying system statistics and quick access to frequently used files.

### Elements

- **Statistics Cards:**

  - **Total Files Card:**

    - Icon: Document
    - Count: 12,486
    - Change: +7.2% from last week (green)

  - **Storage Used Card:**

    - Icon: Cube
    - Amount: 2.34 GB
    - Change: +3.1% from last week (green)

  - **Organization Level Card:**

    - Icon: Checkbox/List
    - Percentage: 72%
    - Change: -15.8% from last week (red)

  - **Smart Managers Card:**
    - Icon: Lock/Manager
    - Count: 3
    - Action: Create new Manager (yellow button)

- **Quick Access Section:**
  - Title: "Quick Access"
  - Files: Study_Guide.pdf (multiple instances)
  - All files modified: 15 May
  - Empty card with "+" icon for adding new quick access item

## 2. Smart Managers Configuration (Image 2)

![Smart Managers Configuration](images/smart_managers.png)

### Purpose

Interface for configuring and managing intelligent file organization rules.

### Elements

- **Header:**

  - Title: "Smart Managers"
  - Action: "Create" button (top right)

- **Manager Configuration Panel:**
  - **Documents Manager:**
    - Status: 75% Organized (green pill)
    - Stats: Managing 324 files across 24 directories
    - Settings:
      - Max number of files per directory: 200 (input field)
      - Max directory depth: 20 (input field)
    - Action Buttons:
      - Sort (button)
      - Rename (button)
      - Delete (button)
      - Save (yellow button)

## 3. Advanced Search (Image 3)

![Advanced Search](images/advanced_search.png)

### Purpose

Comprehensive search interface for finding files using multiple criteria.

### Elements

- **Header:**

  - Title: "Advanced Search"
  - Action: "Clear All" button (top right)

- **Search Interface:**

  - Global search bar: "Search for files, folders, content..."

- **Filter Options:**

  - **File Type:** All Types (dropdown)
  - **Location:** All Locations (dropdown)
  - **Date Modified:** Any Time (dropdown)
  - **Size:** Any Size (dropdown)
  - **Tags:** Any Tags (dropdown)
  - **Author/Owner:** Anyone (dropdown)

- **Results Section:**
  - Title: "Results"
  - Files:
    - Study_Guide.pdf (Modified: 15 May)
    - index.html (Modified: 6 May)

## 4. Settings (Image 4)

![Settings](images/settings.png)

### Purpose

Configuration interface for application preferences and user settings.

### Elements

- **Header:**

  - Title: "Settings"

- **Settings Navigation:**

  - **General** (selected, yellow highlight)
  - Appearance
  - Notifications
  - Account

- **Content Area:**
  - Empty content area (current view shows General tab selected)

## 5. Documents Manager - Folder View (Image 5)

![Documents Manager - Folder View](images/folder_view.png)

### Purpose

Traditional folder-based interface for viewing and managing files within a smart manager.

### Elements

- **Header:**

  - Title: "Documents Manager"
  - Status: 75% Organized (green pill)
  - View Toggle: Folder View (selected), Graph View
  - Action: Sort

- **Search and Filter:**

  - Search: "Search within manager"
  - Sort By: Name (dropdown)
  - Filter: All Files (dropdown)

- **Files and Folders:**

  - Files:
    - Study_Guide.pdf (Modified: 15 May)
    - index.html (Modified: 6 May)
    - notes.docx (Modified: 25 April)
  - Folders:
    - Projects (56 files)
    - Documents (22 files)

- **File Details Panel:**
  - File: notes.docx
  - Type: Word Document
  - Location: Root/Doc
  - Size: 2.8 MB
  - Created: May 10, 2025
  - Modified: May 12, 2025
  - Tags: doc (red), Project (green)
  - Actions: Open File, Add Tags, Lock File, Delete File

## 6. Documents Manager - Graph View (Image 6)

![Documents Manager - Graph View](images/graph_view.png)

### Purpose

Visualization interface showing file relationships in a network graph format.

### Elements

- **Header:**

  - Same as Folder View but with Graph View selected

- **Visualization:**
  - Node graph showing file relationships
  - Central node (yellow) representing focal point or primary file
  - Connected nodes (blue, red, green) representing related files
  - File details panel (same as in Folder View)

## Interaction Notes

1. **Navigation:** Users navigate between screens using the sidebar menu.
2. **Smart Managers:** Each smart manager can be configured with specific rules for organization.
3. **View Modes:** Files can be viewed in both traditional folder view and relationship graph view.
4. **Search Functionality:** Advanced search allows filtering across multiple parameters simultaneously.
5. **Quick Access:** Dashboard provides shortcuts to frequently used files.
6. **File Details:** Sidebar panel displays metadata and provides file actions.
7. **Statistics:** Dashboard provides visual metrics on file management status and trends.

## Image Inclusion Instructions

To include the wireframe images in your document:

1. **Create an images directory**: Create a folder named "images" in the same directory as your markdown file.

2. **Save wireframe images**: Save each wireframe image with a descriptive filename:

   - Image 1: `dashboard.png`
   - Image 2: `smart_managers.png`
   - Image 3: `advanced_search.png`
   - Image 4: `settings.png`
   - Image 5: `folder_view.png`
   - Image 6: `graph_view.png`
   - Sidebar (common to all): `sidebar_navigation.png`

3. **Alternative approach**: If you want to use the actual numbered images as-is:

   ```markdown
   ![Dashboard View](images/image1.png)
   ![Smart Managers Configuration](images/image2.png)
   ![Advanced Search](images/image3.png)
   ![Settings](images/image4.png)
   ![Documents Manager - Folder View](images/image5.png)
   ![Documents Manager - Graph View](images/image6.png)
   ```

4. **For absolute paths**: If your images are stored elsewhere or will be hosted online, use absolute paths:

   ```markdown
   ![Dashboard View](https://yourdomain.com/images/dashboard.png)
   ```

5. **Image sizing**: To control image size, add HTML attributes:
   ```markdown
   <img src="images/dashboard.png" alt="Dashboard View" width="800" />
   ```
