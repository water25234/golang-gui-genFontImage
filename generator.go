package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

func (g *generator) generatorFontImage(pinCodeList string, folder string) (result string, err error) {
	g.Lock()
	defer g.Unlock()
	if len(pinCodeList) == 0 {
		return "", fmt.Errorf("flags readfile is empty")
	}

	if len(folder) == 0 {
		return "", fmt.Errorf("flags folder is empty")
	}

	g.pinCodeList = pinCodeList
	g.folder = folder
	g.fileExt = ".png"

	result = g.processFontImage()
	return result, nil
}

func (g *generator) processFontImage() (result string) {
	fmt.Println("--------------- start work ---------------")

	fileContentArr := strings.Split(g.pinCodeList, "\n")
	fileContentCount := len(fileContentArr)
	errGenCode := &errLog{}

	os.MkdirAll(g.folder, os.ModePerm)

	// channel for job
	jobChans := make(chan jobChannel, fileContentCount)

	// start workers
	wg := &sync.WaitGroup{}
	wg.Add(fileContentCount)

	// start workers
	for i := 1; i <= runtime.NumCPU(); i++ {
		go func(i int) {
			for job := range jobChans {
				g.work(job.fileContent, errGenCode)
				wg.Done()
			}
		}(i)
	}

	// collect job
	for i := 0; i < fileContentCount; i++ {
		jobChans <- jobChannel{
			index:       i,
			fileContent: fileContentArr[i],
		}
	}

	close(jobChans)

	wg.Wait()

	if len(errGenCode.errGenCode) > 0 {
		fmt.Println("error gen font image failure list : ", errGenCode.errGenCode)
	}

	fmt.Println("--------------- finish work ---------------")
	return fmt.Sprintf("執行完成，請找資料夾 『 %s 』 並且確認檔案數量與內容", g.folder)
}

func (g *generator) work(fileContent string, errGenCode *errLog) {

	if len(fileContent) == 0 {
		return
	}

	valueArr := strings.Split(strings.TrimSpace(fileContent), " ")
	valueName, valuePinCode, err := g.pinCodeInfo(valueArr)
	if err != nil {
		return
	}

	pingCode := g.folder + "/" + valueName + g.fileExt

	err = g.fontToImage(
		pingCode,
		valuePinCode,
	)

	if err != nil {
		fmt.Println("gen font image failure", pingCode)
		errGenCode.errGenCode = append(errGenCode.errGenCode, pingCode)
		return
	}

	size, err := g.fileSize(pingCode)
	if err != nil {
		fmt.Println("get file size failure", pingCode)
		errGenCode.errGenCode = append(errGenCode.errGenCode, pingCode)
		return
	}

	fmt.Println(fmt.Sprintf("file: %s, file size: %d", pingCode, size))
	return
}

func (g *generator) pinCodeInfo(valueArr []string) (valueName string, valuePinCode string, err error) {
	if len(valueArr) == 1 {
		valueName = valueArr[0]
		valuePinCode = valueArr[0]
	} else if len(valueArr) == 2 {
		valueName = valueArr[0]
		valuePinCode = valueArr[1]
	} else {
		fmt.Println("value format is error")
		return "", "", fmt.Errorf("value format is error")
	}
	return valueName, valuePinCode, nil
}

func (g *generator) fileSize(pingCode string) (size int64, err error) {
	fi, err := os.Stat(pingCode)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return fi.Size(), nil
}
