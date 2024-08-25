This is the README file for creating a movie database in Golang using the [sqlite](https://gitlab.com/cznic/sqlite) package. This sqlite implementation does require C implementation with Go with cgo, and thus does not suffer performance penalties.

We build a movie database using the publicly available IMDB dataset.


Issues that had to be addressed:
1. Read.ReadAll() encounters parse error on line 355330, column 30: extraneous or missing " in quoted-field. We see that in line 355330, "22-2000" Cidade Aberta", and thus remove the extra paranthesis. We quickly see a limitation with this reader compared to Python's pandas library is the inabllity to handle row errors. Removing this error we are quickly encountered with another error with "Abeltje" (1998/II)",1998,NULL, at line 355609. We quickly see that this is a recurring issue, because the titles with double quotes are not contained within two double quotes. Reader.ReadAll() does not provide any parameters to handle this issue and thus we must rely on writing a custom function or preprocessing the file itself.

The approach is to write a custom function that reads the rows of the csv one by one and inserts into the table, skipping over the rows that result in an error. This insertion of 390,000 rows takes total of 3 minutes and 58 seconds to run, which is notoriously slow. A better approach that is possible may be to preprocess the csv file before hand and use Reader.ReadAll() to read the csv file. However, we are still faced with the problem of inserting the values into table. A possible method is to first parse the CSV into a slice structure and using bulk insert from the Gorm library to do this. 

Possible improvements to this script could be to offer more freedom in user inputs, such as allowing users to write their own SQL scripts, and adding more dimensional tables that can join to the movie table.

To run this, pull the repository, navigate to the directory on the command line and run 'go run main.go'. Enter 1 first to truncate and load the csv file. After loading the table, you can input 2 and find the movie ID you want.