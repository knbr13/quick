// Package quick provides a minimalistic and high-performance web framework for Go.
//
// This file contains example-based tests for various functionalities of the Quick framework,
// including request handling, JSON parsing, XML responses, and header manipulations.
//
// These examples demonstrate how to use Quick's context (`Ctx`) methods effectively
package quick

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// This function is named ExampleCtx_GetReqHeadersAll()
// it with the Examples type.
func ExampleCtx_GetReqHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		// Retrieve all request headers
		headers := c.GetReqHeadersAll()

		// Print specific headers for demonstration
		fmt.Println(headers["Content-Type"]) // Expected: application/json
		fmt.Println(headers["Accept"])       // Expected: application/xml
		return nil
	})

	// Simulate a GET request with headers
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/headers",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/xml"},
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	fmt.Println("Status:", res.StatusCode())

	// Output:
	// [application/json]
	// [application/xml]
	// Status: 200
}

// This function is named ExampleCtx_GetHeadersAll()
// it with the Examples type.
func ExampleCtx_GetHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		// Retrieve all headers from the request
		headers := c.GetHeadersAll()
		fmt.Println(headers["Content-Type"]) // Expected: application/json
		fmt.Println(headers["Accept"])       // Expected: application/xml
		return nil
	})

	// Simulate a GET request with headers
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/headers",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/xml"},
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	// Output:
	// [application/json]
	// [application/xml]
}

// This function is named ExampleCtx_Bind()
// it with the Examples type.
func ExampleCtx_Bind() {
	q := New()

	q.Post("/bind", func(c *Ctx) error {
		// Define a struct to map the JSON request body
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// Parse JSON body into struct
		err := c.Bind(&data)
		if err != nil {
			fmt.Println("Error in Bind:", err)
			return err
		}

		// Print extracted data in desired format
		fmt.Printf("%s %d\n", data.Name, data.Age)
		return nil
	})

	// Simulate a POST request with JSON data
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPost,
		URI:    "/bind",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"name": "Quick", "age": 30}`),
	})

	fmt.Println(res.BodyStr())

	// Output: Quick 30
}

// This function is named ExampleCtx_BodyParser()
// it with the Examples type.
func ExampleCtx_BodyParser() {
	q := New()

	q.Post("/parse", func(c *Ctx) error {
		// Define a struct for JSON parsing
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// Parse request body into the struct
		err := c.BodyParser(&data)
		if err != nil {
			fmt.Println("Erro ao analisar o corpo:", err)
			return err
		}

		// Print parsed data
		fmt.Printf("%s %d\n", data.Name, data.Age)
		return nil
	})

	// Simulate a POST request with JSON data
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPost,
		URI:    "/parse",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"name": "Quick", "age": 28}`),
	})

	fmt.Println(res.BodyStr())

	// Output: Quick 28
}

// This function is named ExampleCtx_Param()
// it with the Examples type.
func ExampleCtx_Param() {
	q := New()

	q.Get("/user/:id", func(c *Ctx) error {
		// Retrieve "id" parameter from the URL path
		id := c.Param("id")
		return c.SendString(id)
	})

	// Simulate a GET request with a path parameter
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/user/42",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: 42
}

// This function is named ExampleCtx_Body()
// it with the Examples type.
func ExampleCtx_Body() {
	q := New()

	// Create a new context with a simulated request body
	q.Post("/body", func(c *Ctx) error {
		// Access the raw body
		body := c.Body()
		fmt.Println(string(body))
		return c.Status(200).String("OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPost,
		URI:    "/body",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"name": "Quick", "age": 28}`),
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	// Output: {"name": "Quick", "age": 28}
}

// This function is named ExampleCtx_Body()
// it with the Examples type.
func ExampleCtx_BodyString() {
	q := New()

	q.Post("/bodyString", func(c *Ctx) error {
		// Access the raw body
		bodyStr := c.BodyString()
		fmt.Println(bodyStr)
		return c.Status(200).String("OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPost,
		URI:    "/bodyString",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"name": "Quick", "age": 28}`),
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	// Output: {"name": "Quick", "age": 28}
}

// This function is named ExampleCtx_JSON()
// it with the Examples type.
func ExampleCtx_JSON() {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		// Define JSON response data
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	// Simulate a GET request and retrieve JSON response
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/json",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err := res.AssertString("{\"message\":\"Hello, Quick!\"}"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: {"message":"Hello, Quick!"}
}

// This function is named ExampleCtx_XML()
// it with the Examples type.
func ExampleCtx_XML() {
	q := New()

	q.Get("/xml", func(c *Ctx) error {
		type XMLResponse struct {
			XMLName xml.Name `xml:"message"`
			Text    string   `xml:",chardata"`
		}

		data := XMLResponse{Text: "Hello, Quick!"}
		xmlBytes, err := xml.Marshal(data)
		if err != nil {
			return c.Status(500).String("Failed to marshal XML")
		}

		c.Set("Content-Type", "application/xml")
		return c.Status(200).Send(xmlBytes)
	})

	// Simulate a GET request and retrieve XML response
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/xml",
		Headers: map[string]string{
			"Accept": "application/xml",
		},
	})

	if err := res.AssertString(`<message>Hello, Quick!</message>`); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: <message>Hello, Quick!</message>
}

// This function is named ExampleCtx_Byte()
// it with the Examples type.
func ExampleCtx_Byte() {
	q := New()

	q.Get("/byte", func(c *Ctx) error {
		// Send raw byte array in response
		return c.Byte([]byte("Hello, Quick!"))
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/byte",
	})

	if err := res.AssertString("Hello, Quick!"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: Hello, Quick!
}

// This function is named ExampleCtx_Send()
// it with the Examples type.
func ExampleCtx_Send() {
	q := New()

	q.Get("/send", func(c *Ctx) error {
		// Send raw bytes in response
		return c.Send([]byte("Hello, Quick!"))
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/send",
	})

	if err := res.AssertString("Hello, Quick!"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: Hello, Quick!
}

// This function is named ExampleCtx_SendString()
// it with the Examples type.
func ExampleCtx_SendString() {
	q := New()

	q.Get("/sendstring", func(c *Ctx) error {
		// Send string response
		return c.SendString("Hello, Quick!")
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/sendstring",
	})

	if err := res.AssertString("Hello, Quick!"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output:	Hello, Quick!
}

// This function is named ExampleCtx_String()
// it with the Examples type.
func ExampleCtx_String() {
	q := New()

	q.Get("/string", func(c *Ctx) error {
		// Return a simple string response
		return c.String("Hello, Quick!")
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/string",
	})

	if err := res.AssertString("Hello, Quick!"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: Hello, Quick!
}

// This function is named ExampleCtx_SendFile()
// it with the Examples type.
func ExampleCtx_SendFile() {
	q := New()

	q.Get("/sendfile", func(c *Ctx) error {
		// Simulate sending a file as a response
		fileContent := []byte("file contents")
		return c.SendFile(fileContent)
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/sendfile",
	})

	if err := res.AssertString("file contents"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: file contents

}

// This function is named ExampleCtx_Set()
// it with the Examples type.
func ExampleCtx_Set() {
	q := New()

	q.Get("/set-header", func(c *Ctx) error {
		// Set a custom response header
		c.Set("X-Custom-Header", "Quick")
		return c.String("Header Set")
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/set-header",
	})

	if err := res.AssertString("Header Set"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.Response().Header.Get("X-Custom-Header"))

	fmt.Println(res.BodyStr())

	// Output:
	// Quick
	// Header Set

}

// This function is named ExampleCtx_Append()
// it with the Examples type.
func ExampleCtx_Append() {
	q := New()

	q.Get("/append-header", func(c *Ctx) error {
		// Append multiple values to a custom header
		c.Append("X-Custom-Header", "Quick")
		return c.String("Header Appended")
	})

	// Simulate a GET request and retrieve XML response
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/append-header",
	})

	if err := res.AssertString("Header Appended"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.Response().Header.Get("X-Custom-Header"))

	fmt.Println(res.BodyStr())

	// Output:
	// Quick
	// Header Appended
}

// This function is named ExampleCtx_Accepts()
// it with the Examples type.
func ExampleCtx_Accepts() {
	q := New()

	q.Get("/accepts", func(c *Ctx) error {
		// Set Accept header
		c.Accepts("application/json")
		return c.String("Accept Set")
	})

	// Simulate a GET request
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/accepts",
	})

	if err := res.AssertString("Accept Set"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.Response().Header.Get("Accept"))

	fmt.Println(res.BodyStr())

	// Output:
	// application/json
	// Accept Set

}

// This function is named ExampleCtx_Status()
// it with the Examples type.
func ExampleCtx_Status() {
	q := New()

	q.Get("/status", func(c *Ctx) error {
		// Set status code to 404
		c.Status(404)
		return c.String("Not Found")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/status", nil, nil)

	fmt.Println(res.Response().StatusCode)

	// Out put: 404
}

// This function is named ExampleCtx_File()
// it with the Examples type.
func ExampleCtx_File() {
	// Creating a Quick instance
	q := New()

	// Defining a route that serves a specific file
	q.Get("/file", func(c *Ctx) error {
		return c.File("quick.txt") // Serves an existing file
	})

	// Simulating a request to test the route
	res, _ := q.QuickTest("GET", "/file", nil)

	// Printing the expected response
	fmt.Println("Status:", res.StatusCode())

	// Out put: Status: 200
}

// This function is named ExampleCtx_JSONIN()
// it with the Examples type.
func ExampleCtx_JSONIN() {
	// Creating a Quick instance
	q := New()

	// Defining a route that returns a structured JSON
	q.Get("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSONIN(data)
	})

	// Simulating a request to test the JSON response
	res, _ := q.QuickTest("GET", "/json", nil)

	// Printing the expected response
	fmt.Println("Status:", res.StatusCode())
	fmt.Println("Body:", res.BodyStr())

	// Out put: Status: 200
	// Body:
	// {
	//  "message": "Hello, Quick!"
	// }
}

// This function is named ExampleCtx_FormFileLimit()
// it with the Examples type.
func ExampleCtx_FormFileLimit() {
	// Creating a new context
	ctx := &Ctx{}

	// Setting a file upload limit to 5MB
	err := ctx.FormFileLimit("5MB")

	// Checking if an error occurred while setting the limit
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Upload limit set to:", ctx.uploadFileSize)
	}

	// Out put: Upload limit set to: 5242880
}

// This function is named ExampleCtx_FormFile()
// it with the Examples type.
func ExampleCtx_FormFile() {
	// Simulated uploaded file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Filename:    "quick.txt",
			Size:        1024,
			ContentType: "text/plain",
			Bytes:       []byte("File content"),
		},
	}

	// Mocking the FormFiles function externally
	mockFormFiles := func(fieldName string) ([]*UploadedFile, error) {
		if fieldName == "file" {
			return []*UploadedFile{uploadedFile}, nil
		}
		return nil, errors.New("file not found")
	}

	// Calling the mocked function instead of modifying `ctx`
	files, err := mockFormFiles("file")

	// Handling the result
	if err != nil {
		fmt.Println("Error:", err)
	} else if len(files) > 0 {
		fmt.Println("Received file:", files[0].FileName())
	}

	// Out put: Received file: quick.txt
}

// This function is named ExampleCtx_FormFiles()
// it with the Examples type.
func ExampleCtx_FormFiles() {
	// Simulating multiple uploaded files
	uploadedFiles := []*UploadedFile{
		{
			Info: FileInfo{
				Filename:    "file1.txt",
				Size:        1024,
				ContentType: "text/plain",
				Bytes:       []byte("File 1 content"),
			},
		},
		{
			Info: FileInfo{
				Filename:    "file2.txt",
				Size:        2048,
				ContentType: "text/plain",
				Bytes:       []byte("File 2 content"),
			},
		},
	}

	// Mocking the FormFiles function externally
	mockFormFiles := func(fieldName string) ([]*UploadedFile, error) {
		if fieldName == "files" {
			return uploadedFiles, nil
		}
		return nil, errors.New("files not found")
	}

	// Calling the mocked function instead of modifying `ctx`
	files, err := mockFormFiles("files")

	// Handling the result
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Received files:")
		for _, file := range files {
			fmt.Printf("- %s (%d bytes)\n", file.FileName(), file.Size())
		}
	}

	// Out put: Received files:
	// - file1.txt (1024 bytes)
	// - file2.txt (2048 bytes)
}

// This function is named ExampleCtx_MultipartForm()
// it with the Examples type.
func ExampleCtx_MultipartForm() {
	// Creating a context and simulating an HTTP request
	ctx := &Ctx{}

	// Simulating an HTTP header with the correct Content-Type
	ctx.Request = &http.Request{
		Header: map[string][]string{
			"Content-Type": {"multipart/form-data"},
		},
	}

	// Attempting to parse the multipart form
	form, err := ctx.MultipartForm()

	// Checking for errors
	if err != nil {
		fmt.Println("Error processing form:", err)
	} else {
		fmt.Println("Form processed successfully:", form)
	}

	// Out put: Form processed successfully: &{...}
}

// This function is named ExampleCtx_GetHeader()
// it with the Examples type.
func ExampleCtx_GetHeader() {
	q := New()

	q.Get("/header", func(c *Ctx) error {
		// Retrieve the "User-Agent" header
		ua := c.GetHeader("User-Agent")
		fmt.Println(ua) // Expected output: "Go-Test-Agent"
		return nil
	})

	// Simulate a GET request with headers
	res, _ := q.QuickTest("GET", "/header", map[string]string{
		"User-Agent": "Go-Test-Agent",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// Go-Test-Agent
}

// This function is named ExampleCtx_GetHeaders()
// it with the Examples type.
func ExampleCtx_GetHeaders() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		// Retrieve all request headers
		headers := c.GetHeaders()

		// Print specific headers for demonstration
		fmt.Println(headers.Get("Content-Type")) // Expected output: "application/json"
		fmt.Println(headers.Get("Accept"))       // Expected output: "application/xml"
		return nil
	})

	// Simulate a GET request with headers
	res, _ := q.QuickTest("GET", "/headers", map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/xml",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// application/json
	// application/xml
}

// This function is named ExampleCtx_RemoteIP()
// it with the Examples type.
func ExampleCtx_RemoteIP() {
	q := New()

	q.Get("/ip", func(c *Ctx) error {
		// Retrieve the client's IP address
		clientIP := c.RemoteIP()

		// Print the IP address for demonstration purposes
		fmt.Println(clientIP)
		return nil
	})

	// Simulate a GET request setting a fixed IP in RemoteAddr
	req := httptest.NewRequest("GET", "/ip", nil)
	req.RemoteAddr = "192.168.1.100:54321" // Setting a fixed IP for testing
	rec := httptest.NewRecorder()

	// Serve the request
	q.ServeHTTP(rec, req)

	// Capture and print the response
	fmt.Println(strings.TrimSpace(rec.Body.String()))

	// Out put:
	// 192.168.1.100
}

// This function is named ExampleCtx_Method()
// it with the Examples type.
func ExampleCtx_Method() {
	q := New()

	q.Post("/method", func(c *Ctx) error {
		fmt.Println(c.Method()) // Expected output: "POST"
		return nil
	})

	// Simulate a POST request
	res, _ := q.QuickTest("POST", "/method", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// POST
}

// This function is named ExampleCtx_Path()
// it with the Examples type.
func ExampleCtx_Path() {
	q := New()

	q.Get("/path/to/resource", func(c *Ctx) error {
		fmt.Println(c.Path()) // Expected output: "/path/to/resource"
		return nil
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/path/to/resource", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// /path/to/resource
}

// This function is named ExampleCtx_QueryParam()
// it with the Examples type.
func ExampleCtx_QueryParam() {
	q := New()

	q.Get("/search", func(c *Ctx) error {
		fmt.Println(c.QueryParam("query")) // Expected output: "quick"
		return nil
	})

	// Simulate a GET request with query parameters
	res, _ := q.QuickTest("GET", "/search?query=quick", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// quick
}
