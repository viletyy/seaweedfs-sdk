package main

import (
	"fmt"
	"os"

	"code.gitlink.org.cn/yystopf/seaweedfs-sdk.git/seaweedfs"
)

func main() {

	filePath := "/Users/virus/admin.sql"

	file, errFile1 := os.Open(filePath)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}

	defer file.Close()

	client := seaweedfs.NewClient("http://localhost:9333", "http://192.168.2.12:8081")

	dirAssign, err := client.GetDirAssign()
	if err != nil {
		fmt.Println(err)
	}

	uploadFile, err := client.PostFile(dirAssign.Fid, file, file.Name())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(uploadFile)
}
