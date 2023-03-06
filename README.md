# git-log-exporter
Export git log to excel file

## Library
github.com/xuri/excelize/v2 ==> write excel file

## How to run
1. Create folder git-repo & export-result inside project directory
2. Checkout git repository that want to read the log, put in git-repo folder
3. Modify variable blow with your environment

  ```go 
  repoBasePathUnix      string = "/d/Project/Learn/git-log-exporter/git-repo/"
  repoBasePathWindows   string = "D:\\Project\\Learn\\git-log-exporter\\git-repo\\"
  resultBasePathUnix    string = "/d/Project/Learn/git-log-exporter/export-result/"
  resultBasePathWindows string = "D:\\Project\\Learn\\git-log-exporter\\export-result\\"
  ```
  
4. Set begin and end date git log that you want to export

  ```go 
  beginDate             string = "2023-02-01"
  endDate               string = "2023-02-28"
  ```
