package main

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const (
	repoBasePathUnix      string = "/d/Project/Learn/git-log-exporter/git-repo/"
	repoBasePathWindows   string = "D:\\Project\\Learn\\git-log-exporter\\git-repo\\"
	resultBasePathUnix    string = "/d/Project/Learn/git-log-exporter/export-result/"
	resultBasePathWindows string = "D:\\Project\\Learn\\git-log-exporter\\export-result\\"
	gitBranch             string = "dlo-dev"
	beginDate             string = "2023-02-01"
	endDate               string = "2023-02-28"
	logFormat             string = "__GIT__SEPARATOR__%x40%h__GIT__DELIMITER__%an__GIT__DELIMITER__%ad__GIT__DELIMITER__%x22%s%x22__GIT__DELIMITER__"
)

var (
	repo = []string{"nbdg-loan-api", "nbdg-loan-kta-api", "nbdg-loan-channeling", "nbdg-channeling-partner", "nbdg-loan-auth", "nbdg-approver-api", "nbdg-loan-data"}
)

type Log struct {
	CommitID     string
	Author       string
	Date         string
	Comment      string
	ChangesFiles string
	LinesAdded   string
	LinesDeleted string
}

func exportExcel(repo string, logs []Log) error {
	f := excelize.NewFile()
	defer func() {
		f.Close()
	}()
	f.SetCellValue("Sheet1", "A1", "Commit ID")
	f.SetCellValue("Sheet1", "B1", "Author")
	f.SetCellValue("Sheet1", "C1", "Date")
	f.SetCellValue("Sheet1", "D1", "Comment")
	f.SetCellValue("Sheet1", "E1", "Changed Files")
	f.SetCellValue("Sheet1", "F1", "Lines Added")
	f.SetCellValue("Sheet1", "G1", "Lines Deleted")

	for i, log := range logs {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), log.CommitID)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), log.Author)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), log.Date)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), log.Comment)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), log.ChangesFiles)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), log.LinesAdded)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), log.LinesDeleted)
	}
	fileName := repo + "_" + beginDate + "_" + endDate + ".xlsx"
	path := resultBasePathUnix
	if runtime.GOOS == "windows" {
		path = resultBasePathWindows
	}
	filePath := path + "/" + fileName

	err := f.SaveAs(filePath)
	if err != nil {
		return err
	}
	return nil
}

func execCommand(binFile string, subcmd string, args ...string) (string, error) {
	arr := append([]string{subcmd}, args...)

	var out bytes.Buffer
	cmd := exec.Command(binFile, arr...)

	//print command string
	//fmt.Println(cmd.String())

	cmd.Stdout = &out
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() != 0 {
				return "", err
			}
		}
	}

	return strings.TrimRight(strings.TrimSpace(out.String()), "\000"), nil
}

func getGitLog(repo string) {
	args := []string{
		fmt.Sprintf(`--since="%s"`, beginDate),
		fmt.Sprintf(`--until="%s"`, endDate),
		"--date=local",
		"--pretty=\"" + logFormat + "\"",
		"--shortstat",
	}

	cwd, err := filepath.Abs(".")

	back := func() error {
		return os.Chdir(cwd)
	}
	defer back()

	if err != nil {
		log.Fatalln(err)
	}

	path := repoBasePathUnix + repo
	if runtime.GOOS == "windows" {
		path = repoBasePathWindows + repo
	}
	dir, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := execCommand("git", "log", args...)
	if err != nil {
		log.Fatalln(err)
	}

	res = sanitizeResult(res)
	segments := strings.Split(res, "__GIT__SEPARATOR__")

	var excelData []Log
	for _, s := range segments {
		properties := strings.Split(s, "__GIT__DELIMITER__")
		if len(properties) < 4 {
			continue
		}
		logData := Log{
			CommitID: properties[0],
			Author:   properties[1],
			Date:     properties[2],
			Comment:  properties[3],
		}
		if len(properties) > 4 && properties[4] != "" {
			properties[4] = strings.ReplaceAll(properties[4], "\" ", "")
			properties[4] = strings.ReplaceAll(properties[4], "\"", "")
			changes := strings.Split(properties[4], ", ")
			if len(changes) > 1 {
				logData.ChangesFiles = changes[0]
				logData.LinesAdded = changes[1]
			}
			if len(changes) > 2 {
				logData.LinesDeleted = changes[2]
			}
		}
		excelData = append(excelData, logData)
	}
	err = exportExcel(repo, excelData)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to export excel. Repo : %s, error : %v", repo, err))
	}
}

func sanitizeResult(res string) string {
	res = strings.ReplaceAll(res, "\n", "")
	res = strings.ReplaceAll(res, "@", "")
	res = strings.ReplaceAll(res, " files changed", "")
	res = strings.ReplaceAll(res, " file changed", "")
	res = strings.ReplaceAll(res, " insertions(+)", "")
	res = strings.ReplaceAll(res, " insertion(+)", "")
	res = strings.ReplaceAll(res, " deletions(-)", "")
	res = strings.ReplaceAll(res, " deletion(-)", "")

	return res
}

func main() {
	for _, r := range repo {
		getGitLog(r)
	}
}
