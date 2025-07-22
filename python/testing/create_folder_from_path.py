import os

# Sample list of paths (can be loaded from a file or generated dynamically)
paths = [    
 "Root/Directory/items/Project_Plan/Apr8TODO.txt",
 "Root/Directory/items/Project_Plan/Apr18 meeting.txt",
 "Root/Directory/items/Project_Plan/COS 301 - Mini-Project - Demo 1 Instructions.pdf ",
 "Root/Directory/items/Project_Plan/COS 301 - Mini-Project - Demo 2 Instructions.pdf ",
 "Root/Directory/items/Project_Plan/COS221 Assignment 1 2025.pdf ",
 "Root/Directory/items/Project_Plan/3.6.4 Survey data to be analysed and visualised for project report mine.xlsx ",
 "Root/Directory/items/Misc/architecture_diagram.png ",
 "Root/Directory/items/Misc/collection_page_wireframe.png ",
 "Root/Directory/items/Misc/login_wireframe.png ",
 "Root/Directory/items/Misc/Screenshot_2025-02-26_at_15.36.48.png ",
 "Root/Directory/items/Misc/statistics_page_wireframe.png ",
 "Root/Directory/items/Misc/UseCase.png ",
 "Root/Directory/items/Misc/~$ecutive summary.docx ",
 "Root/Directory/items/Misc/most challanging.png ",
 "Root/Directory/items/Misc/Taiichi ohno.jpeg ",
 "Root/Directory/items/Program/Assignment2.pdf ",
 "Root/Directory/items/Program/L01_Ch01a(1).pdf ",
 "Root/Directory/items/random_items_inside./COS122 Tutorial 4 Sept 7-8, 2023.pdf ",
 "Root/Directory/items/random_items_inside./~WRL1847.tmp ",
 "Root/Directory/items/Project_Statement.The_garage/User_account_details/cpp_api.md ",
 "Root/Directory/items/Project_Statement.The_garage/User_account_details/mp11_design_specification.md ",
 "Root/Directory/items/Project_Statement.The_garage/User_account_details/mp11_requirement_spec.md ",
 "Root/Directory/items/Project_Statement.The_garage/User_account_details/TODO mar30 Meeting.txt ",
 "Root/Directory/items/Project_Statement.The_garage/User_account_details/Tututorial_2.pdf ",
 "Root/Directory/items/Project_Statement.The_garage/~WRL0005.tmp ",
 "Root/Directory/items/Project_Statement.The_garage/Presentation speech.docx ",
 "Root/Directory/items/Project_Statement.The_garage/Project Budget Form 2024.pdf ",
 "Root/Directory/items/distribution/Week 3_Tutorial_2024_with Answers.pdf ",
 "Root/Directory/items/distribution/Week 4_Tutorial_with answers.pdf ",
 "Root/Directory/items/distribution/Week 5_Tutorial_2024_with answers.pdf ",
 "Root/Directory/items/JSON_object_arrays/Misc/DeeBee.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Gantt chart.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Most rewarding.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Picture2.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Document[1].pdf ",
 "Root/Directory/items/JSON_object_arrays/Misc/ENjoyment.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Gauteng.png ",
 "Root/Directory/items/JSON_object_arrays/Misc/Picture1.png ",
 "Root/Directory/items/JSON_object_arrays/Importing the Database.md ",
 "Root/Directory/items/JSON_object_arrays/L05_Ch02c.pdf ",
 "Root/Directory/items/JSON_object_arrays/MP Progress report.txt ",
 "Root/Directory/items/JSON_object_arrays/MPChecklist.txt ",
 "Root/Directory/items/JSON_object_arrays/Prac1Triggers.txt"
]

# Create folder structure and dummy files
for path in paths:
    dir_path = os.path.dirname(path)
    os.makedirs(dir_path, exist_ok=True)

    # Create dummy file with placeholder content
    with open(path, 'w') as f:
        f.write(f"Dummy content for {os.path.basename(path)}")

print("All folders and files created.")
