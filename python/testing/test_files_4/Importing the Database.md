Importing the Database Dump

To import the database dump into MariaDB, follow these steps:

1. Open the terminal or command prompt.
2. Run the following command to create a new database:
   
   CREATE DATABASE u23536030_carhiresystem;

3. Import the database using `mysqldump`:

   mysql -u root -p u23536030_carhiresystem < u23536030_carhiresystem_dump.sql

4. Enter your MariaDB root password when prompted. (COS221)

The database is now restored successfully.

