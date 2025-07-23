import os

# Sample list of paths (can be loaded from a file or generated dynamically)
paths = [    
"Root2/Directory/code_program_task/project_plan_design/Apr8TODO.txt",
"Root2/Directory/code_program_task/project_plan_design/Apr18 meeting.txt",
"Root2/Directory/code_program_task/project_plan_design/COS 301 - Mini-Project - Demo 1 Instructions.pdf",
"Root2/Directory/code_program_task/project_plan_design/COS 301 - Mini-Project - Demo 2 Instructions.pdf",
"Root2/Directory/code_program_task/project_plan_design/COS221 Assignment 1 2025.pdf",
"Root2/Directory/code_program_task/project_plan_design/3.6.4 Survey data to be analysed and visualised for project report mine.xlsx",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/architecture_diagram.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/collection_page_wireframe.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/login_wireframe.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/Screenshot_2025-02-26_at_15.36.48.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/statistics_page_wireframe.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/UseCase.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/~$ecutive summary.docx",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/most challanging.png",
"Root2/Directory/code_program_task/statistics_page_wireframe_usecase/Taiichi ohno.jpeg",
"Root2/Directory/code_program_task/program_code_assignment/Assignment2.pdf",
"Root2/Directory/code_program_task/program_code_assignment/L01_Ch01a(1).pdf",
"Root2/Directory/code_program_task/process_question_data/COS122 Tutorial 4 Sept 7-8, 2023.pdf",
"Root2/Directory/code_program_task/process_question_data/~WRL1847.tmp",
"Root2/Directory/code_program_task/item_data_project/item_data_query/query_database/cpp_api.md",
"Root2/Directory/code_program_task/item_data_project/item_data_query/query_database/mp11_design_specification.md",
"Root2/Directory/code_program_task/item_data_project/item_data_query/query_database/mp11_requirement_spec.md",
"Root2/Directory/code_program_task/item_data_project/item_data_query/query_database/TODO mar30 Meeting.txt",
"Root2/Directory/code_program_task/item_data_project/item_data_query/query_database/Tututorial_2.pdf",
"Root2/Directory/code_program_task/item_data_project/item_data_query/item_objectivesthe_project/~WRL0005.tmp",
"Root2/Directory/code_program_task/item_data_project/item_data_query/item_objectivesthe_project/Presentation speech.docx",
"Root2/Directory/code_program_task/item_data_project/item_data_query/item_objectivesthe_project/Project Budget Form 2024.pdf",
"Root2/Directory/code_program_task/item_data_project/probability_event_distribution/Week 3_Tutorial_2024_with Answers.pdf",
"Root2/Directory/code_program_task/item_data_project/probability_event_distribution/Week 4_Tutorial_with answers.pdf",
"Root2/Directory/code_program_task/item_data_project/probability_event_distribution/Week 5_Tutorial_2024_with answers.pdf",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture2_most_rewarding/gantt_chart_deebee/DeeBee.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture2_most_rewarding/gantt_chart_deebee/Gantt chart.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture2_most_rewarding/most_rewarding_picture2/Most rewarding.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture2_most_rewarding/most_rewarding_picture2/Picture2.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture1_enjoyment/Document[1].pdf",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture1_enjoyment/ENjoyment.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture1_enjoyment/Gauteng.png",
"Root2/Directory/code_program_task/operating_system/picture2_picture1/picture1_enjoyment/Picture1.png",
"Root2/Directory/code_program_task/operating_system/operating_system/Importing the Database.md",
"Root2/Directory/code_program_task/operating_system/operating_system/L05_Ch02c.pdf",
"Root2/Directory/code_program_task/operating_system/operating_system/MP Progress report.txt",
"Root2/Directory/code_program_task/operating_system/operating_system/MPChecklist.txt",
"Root2/Directory/code_program_task/operating_system/operating_system/Prac1Triggers.txt"
]

# Create folder structure and dummy files
for path in paths:
    dir_path = os.path.dirname(path)
    os.makedirs(dir_path, exist_ok=True)

    # Create dummy file with placeholder content
    with open(path, 'w') as f:
        f.write(f"Dummy content for {os.path.basename(path)}")

print("All folders and files created.")
