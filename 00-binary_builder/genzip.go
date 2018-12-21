package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func run(name string, args ...string) error {
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")

	cmd := exec.Command(name, args...)
	log.Printf("Waiting for command to finish...\n")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Run: Command Execution Error: ", err)
	}
	log.Printf("Done!\n")
	return nil
}

func zip(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Zip: Open Error: ", err)
	}
	defer f.Close()
	ds, err := os.OpenFile(path+".zip", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0466)
	if err != nil {
		return err
	}
	defer ds.Close()
	zf := zlib.NewWriter(ds)
	defer zf.Close()
	if _, err = io.Copy(zf, f); err != nil {
		return err
	}
	return nil
}

func main() {
	fpath := "main"
	if len(os.Args) > 1 {
		fpath = os.Args[1]
	}
	fmt.Println("Building: " + fpath + ".go")
	err := run("go", "build", "-o", fpath, fpath+".go")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Zipping: " + fpath + ".zip")
	err = zip(fpath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Binary built with success.")
}
