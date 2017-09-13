package controllers

import (
	"github.com/revel/revel"
	"store-read-image/app/database"
	"io"
	"os"
	"fmt"
	"bytes"
)

type Gimage struct {
	
	*revel.Controller
}

func (c Gimage) Upload() revel.Result{
	
	
	//Read image from Request body
	c.Request.ParseMultipartForm(32 << 20)
	file, handler, err := c.Request.FormFile("uploadfile")
    check(err)
    defer file.Close()
	
	//Convert Multipart Form File to byte
	
	buf := bytes.NewBuffer(nil)
	io.Copy(buf,file)

	
	//Store Image File in MongoDB GridFS
    f, err := database.Gimage.Create(handler.Filename)
    check(err)
	n ,err := f.Write(buf.Bytes())
	check (err)
    defer f.Close()
	fmt.Println("Written: ", n , "Bytes")
	return c.RenderJSON("File stored in mongo db")
	
}

func (c Gimage) Read(imagename string) revel.Result{
	
	//Get the size of the image
	file, err := database.Gimage.Open(imagename)
	check(err)
	b := make([]byte,file.Size())
	
	//Read Image Contents in bytes
	n,err := file.Read(b)
	check(err)
	fmt.Println(n)
	defer file.Close()
	
	//Write to an image file
	f, err := os.OpenFile("./store-read-image/"+imagename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	_,err = f.Write(b)
	check(err)
	defer f.Close()
	
	return c.RenderJSON("image file writtern")
}


func check(err error){
	
	if err != nil {
		
		fmt.Println(err)
		
	}
}