import os

# Sample list of paths (can be loaded from a file or generated dynamically)
paths = [    
"test_files_3/code_python/project_plan/Apr8TODO.txt",
"test_files_3/code_python/project_plan/Apr18 meeting.txt",
"test_files_3/code_python/project_plan/COS 301 - Mini-Project - Demo 1 Instructions.pdf",
"test_files_3/code_python/project_plan/COS 301 - Mini-Project - Demo 2 Instructions.pdf",
"test_files_3/code_python/project_plan/COS221 Assignment 1 2025.pdf",
"test_files_3/code_python/project_plan/3.6.4 Survey data to be analysed and visualised for project report mine.xlsx",
"test_files_3/code_python/python_statistics_page_wireframe/architecture_diagram.png",
"test_files_3/code_python/python_statistics_page_wireframe/collection_page_wireframe.png",
"test_files_3/code_python/python_statistics_page_wireframe/login_wireframe.png",
"test_files_3/code_python/python_statistics_page_wireframe/Screenshot_2025-02-26_at_15.36.48.png",
"test_files_3/code_python/python_statistics_page_wireframe/statistics_page_wireframe.png",
"test_files_3/code_python/python_statistics_page_wireframe/UseCase.png",
"test_files_3/code_python/python_statistics_page_wireframe/~$ecutive summary.docx",
"test_files_3/code_python/python_statistics_page_wireframe/most challanging.png",
"test_files_3/code_python/python_statistics_page_wireframe/Taiichi ohno.jpeg",
"test_files_3/code_python/code_program/Assignment2.pdf",
"test_files_3/code_python/code_program/L01_Ch01a(1).pdf",
"test_files_3/code_python/question_process/COS122 Tutorial 4 Sept 7-8, 2023.pdf",
"test_files_3/code_python/question_process/~WRL1847.tmp",
"test_files_3/code_python/data_item/cpp_api.md",
"test_files_3/code_python/data_item/mp11_design_specification.md",
"test_files_3/code_python/data_item/mp11_requirement_spec.md",
"test_files_3/code_python/data_item/TODO mar30 Meeting.txt",
"test_files_3/code_python/data_item/Tututorial_2.pdf",
"test_files_3/code_python/data_item/~WRL0005.tmp",
"test_files_3/code_python/data_item/Presentation speech.docx",
"test_files_3/code_python/data_item/Project Budget Form 2024.pdf",
"test_files_3/code_python/data_item/probability_event/Week 3_Tutorial_2024_with Answers.pdf",
"test_files_3/code_python/data_item/probability_event/Week 4_Tutorial_with answers.pdf",
"test_files_3/code_python/data_item/probability_event/Week 5_Tutorial_2024_with answers.pdf",
"test_files_3/code_python/operating_system/python_picture2/picture2_most/gantt_chart/DeeBee.png",
"test_files_3/code_python/operating_system/python_picture2/picture2_most/gantt_chart/Gantt chart.png",
"test_files_3/code_python/operating_system/python_picture2/picture2_most/most_rewarding/Most rewarding.png",
"test_files_3/code_python/operating_system/python_picture2/picture2_most/most_rewarding/Picture2.png",
"test_files_3/code_python/operating_system/python_picture2/picture1_enjoyment/Document[1].pdf",
"test_files_3/code_python/operating_system/python_picture2/picture1_enjoyment/ENjoyment.png",
"test_files_3/code_python/operating_system/python_picture2/picture1_enjoyment/Gauteng.png",
"test_files_3/code_python/operating_system/python_picture2/picture1_enjoyment/Picture1.png",
"test_files_3/code_python/operating_system/Importing the Database.md",
"test_files_3/code_python/operating_system/L05_Ch02c.pdf",
"test_files_3/code_python/operating_system/MP Progress report.txt",
"test_files_3/code_python/operating_system/MPChecklist.txt",
"test_files_3/code_python/operating_system/Prac1Triggers.txt"
]

# Create folder structure and dummy files
for path in paths:
    dir_path = os.path.dirname(path)
    os.makedirs(dir_path, exist_ok=True)

    # Create dummy file with placeholder content
    with open(path, 'w') as f:
        f.write(f"Dummy content for {os.path.basename(path)}")

print("All folders and files created.")
