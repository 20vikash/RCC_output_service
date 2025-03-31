package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	o "outputservice/grpc/server"
)

func (app *Application) OutputCCpp(ctx context.Context, req *o.OutputRequest) (*o.OutputResponse, error) {
	var cmd *exec.Cmd
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	var res string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.Language == "cpp" {
		req.Language = "c++"
	}

	fileName := req.RoomID + req.UserName

	cmd = exec.Command("g++", "-x", req.Language, "-", "-o", fileName)
	cmd.Stdin = bytes.NewBufferString(req.Code)
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		res = stdErr.String()
		return &o.OutputResponse{Message: res}, nil
	}

	defer os.Remove(fileName)

	cmd = exec.CommandContext(ctx, "./"+fileName)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()

	if err != nil {
		if err == err.(*exec.ExitError) {
			res = "Took too long to generate the output"
			return &o.OutputResponse{Message: res}, nil
		}
		log.Println(err.Error())
		res = stdErr.String()
		return &o.OutputResponse{Message: res}, nil
	} else {
		res = stdOut.String()
		return &o.OutputResponse{Message: res}, nil
	}
}

func (app *Application) OutputPython(ctx context.Context, req *o.OutputRequest) (*o.OutputResponse, error) {
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3")
	cmd.Stdin = bytes.NewBufferString(req.Code)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		if err == err.(*exec.ExitError) && stdErr.String() == "" {
			return &o.OutputResponse{Message: "Took too long to generate the output"}, nil
		}
		return &o.OutputResponse{Message: stdErr.String()}, nil
	}

	return &o.OutputResponse{Message: stdOut.String()}, nil
}

func (app *Application) OutputGolangPHP(ctx context.Context, req *o.OutputRequest) (*o.OutputResponse, error) {
	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	var extension string
	var cmd *exec.Cmd
	var file string
	var filePath string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	file = req.RoomID + req.UserName

	if req.Language == "go" {
		extension = ".go"
		filePath = "../sandbox/" + file + extension
		cmd = exec.CommandContext(ctx, "go", "run", filePath)
	} else {
		extension = ".php"
		filePath = "../sandbox/" + file + extension
		cmd = exec.CommandContext(ctx, "php", filePath)
	}

	if err := os.WriteFile(filePath, []byte(req.Code), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return &o.OutputResponse{Message: err.Error()}, nil
	}

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running file:", err)

		if err := os.Remove(filePath); err != nil {
			fmt.Println("Error removing file:", err)
		}

		return &o.OutputResponse{Message: stdErr.String()}, nil
	}

	if err := os.Remove(filePath); err != nil {
		fmt.Println("Error removing file:", err)
	}

	return &o.OutputResponse{Message: stdOut.String()}, nil
}
